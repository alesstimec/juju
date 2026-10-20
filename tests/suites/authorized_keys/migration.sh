run_migrate_authorized_keys() {
	SUB_TEST_DIR="${TEST_DIR}/migrate_authorized_keys"
	mkdir -p "${SUB_TEST_DIR}"
	log_file="${SUB_TEST_DIR}/migrate-keys.log"

	# record the controller made before this test
	to_controller=$BOOTSTRAPPED_JUJU_CTRL_NAME

	BOOTSTRAP_REUSE=false \
	bootstrap "migrate-keys" "$log_file"
	from_controller=$BOOTSTRAPPED_JUJU_CTRL_NAME

	key_file="${SUB_TEST_DIR}/id_ed25519"
	key_file_pub="${key_file}.pub"
	ssh-keygen -t ed25519 -f "$key_file" -C "isgreat@juju.is" -P ""
	fingerprint=$(ssh-keygen -lf "${key_file_pub}" | cut -f 2 -d ' ')

	juju add-ssh-key "$(cat ${key_file_pub})"
	check_contains "$(juju ssh-keys)" "$fingerprint"

	# Wait for the bootstrap machine's instance status to be reported as
	# running before migrating. The instance poller only confirms the
	# instance status a short time after bootstrap, and the migration
	# source precheck rejects controllers whose machines do not have
	# instance status running.
	wait_for_instance_status "0" "running" "${from_controller}:controller"

	juju migrate "migrate-keys" "${to_controller}"

	juju switch "${to_controller}:controller"

	## Wait for the new model migration to appear in the to controller.
	wait_for_model "migrate-keys"

	juju switch "${to_controller}:migrate-keys"
	check_contains "$(juju ssh-keys)" "$fingerprint"

	# The migrated machine must carry the target controller's own key so
	# that the target controller can reach it directly (for example during
	# disaster recovery), and must no longer carry the source controller's.
	# The machine agent reconciles the key shortly after it reconnects to
	# the target controller, so poll rather than check once.
	wait_for_machine_agent_status "0" "started"
	target_system_key=$(juju controller-config -c "${to_controller}" system-ssh-keys | awk '{print $1" "$2}')
	source_system_key=$(juju controller-config -c "${from_controller}" system-ssh-keys | awk '{print $1" "$2}')
	check_migrated_machine_controller_key() {
		local machine_keys
		# shellcheck disable=SC2016
		machine_keys=$(juju ssh -m "${to_controller}:migrate-keys" 0 'cat ~/.ssh/authorized_keys')
		check_contains "${machine_keys}" "${target_system_key}" &&
			check_not_contains "${machine_keys}" "${source_system_key}"
	}
	retry check_migrated_machine_controller_key 12 5

	destroy_controller "$from_controller"
}

test_migrate_authorized_keys() {
	if [ "$(skip 'test_migrate_authorized_keys')" ]; then
		echo "==> TEST SKIPPED: migrate authorized keys"
		return
	fi

	(
		set_verbosity

		run "run_migrate_authorized_keys"
	)
}
