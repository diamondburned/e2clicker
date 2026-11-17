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

with subtest("postgresql schema version and seeded rows exist"):
    meta_v = machine.succeed(
        'runuser -u postgres -- psql -d e2clicker-backend -t -A -c "SELECT v FROM meta;"'
    ).strip()
    assert meta_v == "4", f"expected meta.v == 4, got {meta_v!r}"

    user_count = int(
        machine.succeed(
            'runuser -u postgres -- psql -d e2clicker-backend -t -A -c "SELECT count(*) FROM users;"'
        ).strip()
    )
    assert user_count >= 3, f"expected >= 3 users from hurl tests, got {user_count}"
