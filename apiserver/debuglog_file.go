// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/juju/juju/state"
	"github.com/juju/loggo"
	"github.com/juju/names"
	"github.com/juju/utils/tailer"
)

func newDebugLogFileHandler(ctxt httpContext, stop <-chan struct{}, logDir string) http.Handler {
	fileHandler := &debugLogFileHandler{logDir: logDir}
	return newDebugLogHandler(ctxt, stop, fileHandler.handle)
}

// debugLogFileHandler handles requests to watch all-machines.log.
type debugLogFileHandler struct {
	logDir string
}

func (h *debugLogFileHandler) handle(
	_ state.LoggingState,
	params *debugLogParams,
	socket debugLogSocket,
	stop <-chan struct{},
) error {
	stream := newLogFileStream(params)

	// Open log file.
	logLocation := filepath.Join(h.logDir, "all-machines.log")
	logFile, err := os.Open(logLocation)
	if err != nil {
		socket.sendError(fmt.Errorf("cannot open log file: %v", err))
		return err
	}
	defer logFile.Close()

	if err := stream.positionLogFile(logFile); err != nil {
		socket.sendError(fmt.Errorf("cannot position log file: %v", err))
		return err
	}

	// If we get to here, no more errors to report.
	if err := socket.sendOk(); err != nil {
		return err
	}

	stream.start(logFile, socket)
	return stream.wait(stop)
}

func newLogFileStream(params *debugLogParams) *logFileStream {
	return &logFileStream{
		debugLogParams:  params,
		maxLinesReached: make(chan bool),
	}
}

type logFileLine struct {
	line      string
	agentTag  string
	agentName string
	level     loggo.Level
	module    string
}

func parseLogLine(line string) *logFileLine {
	const (
		agentTagIndex = 0
		levelIndex    = 3
		moduleIndex   = 4
	)
	fields := strings.Fields(line)
	result := &logFileLine{
		line: line,
	}
	if len(fields) > agentTagIndex {
		agentTag := fields[agentTagIndex]
		// Drop mandatory trailing colon (:).
		// Since colon is mandatory, agentTag without it is invalid and will be empty ("").
		if strings.HasSuffix(agentTag, ":") {
			result.agentTag = agentTag[:len(agentTag)-1]
		}
		/*
		 Drop unit suffix.
		 In logs, unit information may be prefixed with either a unit_tag by itself or a unit_tag[nnnn].
		 The code below caters for both scenarios.
		*/
		if bracketIndex := strings.Index(agentTag, "["); bracketIndex != -1 {
			result.agentTag = agentTag[:bracketIndex]
		}
		// If, at this stage, result.agentTag is empty,  we could not deduce the tag. No point getting the name...
		if result.agentTag != "" {
			// Entity Name deduced from entity tag
			entityTag, err := names.ParseTag(result.agentTag)
			if err != nil {
				/*
				 Logging error but effectively swallowing it as there is no where to propogate.
				 We don't expect ParseTag to fail since the tag was generated by juju in the first place.
				*/
				logger.Errorf("Could not deduce name from tag %q: %v\n", result.agentTag, err)
			}
			result.agentName = entityTag.Id()
		}
	}
	if len(fields) > moduleIndex {
		if level, valid := loggo.ParseLevel(fields[levelIndex]); valid {
			result.level = level
			result.module = fields[moduleIndex]
		}
	}

	return result
}

// logFileStream runs the tailer to read a log file and stream it via
// a web socket.
type logFileStream struct {
	*debugLogParams
	logTailer       *tailer.Tailer
	lineCount       uint
	maxLinesReached chan bool
}

// positionLogFile will update the internal read position of the logFile to be
// at the end of the file or somewhere in the middle if backlog has been specified.
func (stream *logFileStream) positionLogFile(logFile io.ReadSeeker) error {
	// Seek to the end, or lines back from the end if we need to.
	if !stream.fromTheStart {
		return tailer.SeekLastLines(logFile, stream.backlog, stream.filterLine)
	}
	return nil
}

// start the tailer listening to the logFile, and sending the matching
// lines to the writer.
func (stream *logFileStream) start(logFile io.ReadSeeker, writer io.Writer) {
	stream.logTailer = tailer.NewTailer(logFile, writer, stream.countedFilterLine)
}

// wait blocks until the logTailer is done or the maximum line count
// has been reached or the stop channel is closed.
func (stream *logFileStream) wait(stop <-chan struct{}) error {
	select {
	case <-stream.logTailer.Dead():
		return stream.logTailer.Err()
	case <-stream.maxLinesReached:
		stream.logTailer.Stop()
	case <-stop:
		stream.logTailer.Stop()
	}
	return nil
}

// filterLine checks the received line for one of the configured tags.
func (stream *logFileStream) filterLine(line []byte) bool {
	log := parseLogLine(string(line))
	return stream.checkIncludeEntity(log) &&
		stream.checkIncludeModule(log) &&
		!stream.exclude(log) &&
		stream.checkLevel(log)
}

// countedFilterLine checks the received line for one of the configured tags,
// and also checks to make sure the stream doesn't send more than the
// specified number of lines.
func (stream *logFileStream) countedFilterLine(line []byte) bool {
	result := stream.filterLine(line)
	if result && stream.maxLines > 0 {
		stream.lineCount++
		result = stream.lineCount <= stream.maxLines
		if stream.lineCount == stream.maxLines {
			close(stream.maxLinesReached)
		}
	}
	return result
}

func (stream *logFileStream) checkIncludeEntity(line *logFileLine) bool {
	if len(stream.includeEntity) == 0 {
		return true
	}
	for _, value := range stream.includeEntity {
		if agentMatchesFilter(line, value) {
			return true
		}
	}
	return false
}

// agentMatchesFilter checks if agentTag tag or agentTag name match given filter
func agentMatchesFilter(line *logFileLine, aFilter string) bool {
	return hasMatch(line.agentName, aFilter) || hasMatch(line.agentTag, aFilter)
}

// hasMatch determines if value contains filter using regular expressions.
// All wildcard occurrences are changed to `.*`
// Currently, all match exceptions are logged and not propagated.
func hasMatch(value, aFilter string) bool {
	/* Special handling: out of 12 regexp metacharacters \^$.|?+()[*{
	   only asterix (*) can be legally used as a wildcard in this context.
	   Both machine and unit tag and name specifications do not allow any other metas.
	   Consequently, if aFilter contains wildcard (*), do not escape it -
	   transform it into a regexp "any character(s)" sequence.
	*/
	aFilter = strings.Replace(aFilter, "*", `.*`, -1)
	matches, err := regexp.MatchString("^"+aFilter+"$", value)
	if err != nil {
		// logging errors here... but really should they be swallowed?
		logger.Errorf("\nCould not match filter %q and regular expression %q\n.%v\n", value, aFilter, err)
	}
	return matches
}

func (stream *logFileStream) checkIncludeModule(line *logFileLine) bool {
	if len(stream.includeModule) == 0 {
		return true
	}
	for _, value := range stream.includeModule {
		if strings.HasPrefix(line.module, value) {
			return true
		}
	}
	return false
}

func (stream *logFileStream) exclude(line *logFileLine) bool {
	for _, value := range stream.excludeEntity {
		if agentMatchesFilter(line, value) {
			return true
		}
	}
	for _, value := range stream.excludeModule {
		if strings.HasPrefix(line.module, value) {
			return true
		}
	}
	return false
}

func (stream *logFileStream) checkLevel(line *logFileLine) bool {
	return line.level >= stream.filterLevel
}
