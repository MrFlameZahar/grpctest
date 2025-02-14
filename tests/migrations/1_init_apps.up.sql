INSERT INTO apps (id, name, secret)
values (1, 'test', 'test-secret')
ON CONFLICT DO NOTHING;