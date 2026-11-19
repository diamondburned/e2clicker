start_all()

with subtest("services start and systemd hardening is active"):
    machine.wait_for_unit("postgresql.service")
    machine.wait_for_unit("e2clicker-backend.service")
    machine.wait_for_unit("e2clicker-frontend.socket")
    machine.wait_for_unit("caddy.service")
    machine.wait_for_open_port(36001)
    machine.wait_for_open_port(80)

    assert (
        machine.succeed("systemctl show e2clicker-backend.service -p DynamicUser --value").strip()
        == "yes"
    )
    assert (
        machine.succeed("systemctl show e2clicker-frontend.service -p DynamicUser --value").strip()
        == "yes"
    )

with subtest("frontend SSR responds through caddy"):
    login_html = machine.succeed("curl -fsSL http://127.0.0.1/login")
    assert "e2clicker" in login_html.lower(), f"unexpected login HTML: {login_html[:500]}"

with subtest("hurl black-box API suite passes through caddy"):
    machine.succeed("hurl --test --variable base_url=http://127.0.0.1 /etc/e2clicker-tests/api/*.hurl")

with subtest("postgresql schema version, index, and reminder attempt semantics"):
    meta_v = machine.succeed(
        'runuser -u postgres -- psql -d e2clicker-backend -t -A -c "SELECT v FROM meta;"'
    ).strip()
    assert meta_v == "7", f"expected meta.v == 7, got {meta_v!r}"

    user_count = int(
        machine.succeed(
            'runuser -u postgres -- psql -d e2clicker-backend -t -A -c "SELECT count(*) FROM users;"'
        ).strip()
    )
    assert user_count >= 3, f"expected >= 3 users from hurl tests, got {user_count}"

    has_idx = machine.succeed(
        "runuser -u postgres -- psql -d e2clicker-backend -t -A -c "
        "\"SELECT count(*) FROM pg_indexes WHERE indexname IN ('notification_history_user_sent_at', 'dosage_history_user_taken_at_id', 'dosage_history_pkey');\""
    ).strip()
    assert has_idx == "3", f"expected 3 indexes on notification_history/dosage_history, got {has_idx!r}"

    old_objs = machine.succeed(
        "runuser -u postgres -- psql -d e2clicker-backend -t -A -c "
        "\"SELECT count(*) FROM pg_indexes WHERE indexname IN ('dosage_history_user_secret_taken_at_key', 'dosage_history_taken_at', 'dosage_history_user_secret');\""
    ).strip()
    assert old_objs == "0", f"expected old dosage_history indexes/constraints dropped, got {old_objs!r}"

    # Verify errored notification attempts are counted as genuine attempts.
    machine.succeed(
        "runuser -u postgres -- psql -d e2clicker-backend -v ON_ERROR_STOP=1 <<'SQL'\n"
        "INSERT INTO users (secret, name) VALUES ('vm-retry-test', 'Retry Test');\n"
        "INSERT INTO dosage_schedule (user_secret, delivery_method, dose, interval)\n"
        "  VALUES ('vm-retry-test', 'EV im', 5, '1 day');\n"
        "INSERT INTO dosage_history (id, user_secret, delivery_method, dose, taken_at)\n"
        "  VALUES ('019450a0-0000-7000-8000-000000000099', 'vm-retry-test', 'EV im', 5, now() - interval '2 days');\n"
        "INSERT INTO notification_history (user_secret, sent_at, supposed_entity_time, error_reason)\n"
        "  VALUES ('vm-retry-test', now() - interval '1 minute', now() - interval '2 days', 'transient 503');\n"
        "SQL"
    )
    reminded_at = machine.succeed(
        "runuser -u postgres -- psql -d e2clicker-backend -t -A -c \""
        "SELECT COALESCE(last_notification.sent_at::text, 'NULL') "
        "FROM users "
        "INNER JOIN dosage_schedule ON users.secret = dosage_schedule.user_secret "
        "INNER JOIN LATERAL ("
        "  SELECT * FROM dosage_history dh WHERE dh.user_secret = users.secret ORDER BY dh.taken_at DESC, dh.id DESC LIMIT 1"
        ") dh ON TRUE "
        "LEFT JOIN LATERAL ("
        "  SELECT nh.sent_at FROM notification_history nh "
        "  WHERE nh.user_secret = users.secret "
        "  ORDER BY nh.supposed_entity_time DESC NULLS LAST, nh.sent_at DESC LIMIT 1"
        ") last_notification ON TRUE "
        "WHERE users.secret = 'vm-retry-test';\""
    ).strip()
    assert reminded_at != "NULL", f"expected errored attempt to be counted, got {reminded_at!r}"
