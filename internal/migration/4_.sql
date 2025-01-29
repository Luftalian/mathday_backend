-- +goose Up

-- イベントデータの挿入
INSERT INTO events (
    title,
    organizer,
    start_date,
    start_time,
    end_date,
    end_time,
    email,
    prefecture,
    event_type,
    is_online,
    is_offline,
    official_url,
    online_lecture_url,
    venue,
    description,
    tags,
    speakers,
    schedule,
    auth_code,
    is_authenticated
)
VALUES
-- 1件目: 数学ミニフォーラム（オフライン）
(
    '数学ミニフォーラム',
    '東京大学 数理教室',
    '2024-03-14',            -- start_date
    '13:00:00',             -- start_time
    '2024-03-14',            -- end_date
    '17:00:00',             -- end_time
    'info@mathforum.example.com', -- email（サンプル）
    '東京都',                -- prefecture
    'オフライン',            -- event_type
    FALSE,                  -- is_online
    TRUE,                   -- is_offline
    'https://example.com/math-forum', -- official_url
    NULL,                   -- online_lecture_url
    '東京都',               -- venue
    '数学の最新トピックについて、専門家が分かりやすく解説します。', -- description
    '["Business","Vacation"]',                   -- tags
    '[{"name":"Heidi","title":"CEO","organization":"GlobalCorp"}]',                   -- speakers
    '[{"time":"19:00","title":"Welcome Speech","speaker":"Heidi"}]',                   -- schedule
    '11111111-1111-1111-1111-111111111111', -- auth_code（例）
    TRUE                    -- is_authenticated
),
-- 2件目: 数学の歴史とパイのお話会（オンライン）
(
    '数学の歴史とパイのお話会',
    '東京大学カブリIPMU',
    '2024-03-14',            -- start_date
    '10:00:00',             -- start_time
    '2024-03-14',            -- end_date
    '12:00:00',             -- end_time
    'history@mathforum.example.com', -- email（サンプル）
    NULL,                   -- prefecture（オンラインのみなのでNULL）
    'オンライン',           -- event_type
    TRUE,                   -- is_online
    FALSE,                  -- is_offline
    'https://example.com/math-history', -- official_url
    'https://example.com/math-history/online', -- online_lecture_url（必要に応じて追加）
    'オンライン',          -- venue（オンラインと記入）
    '数学の歴史を紐解きながら、円周率πの魅力に迫ります。', -- description
    '["Business","Vacation"]',                   -- tags
    '[{"name":"Heidi","title":"CEO","organization":"GlobalCorp"}]',                   -- speakers
    '[{"time":"19:00","title":"Welcome Speech","speaker":"Heidi"}]',                   -- schedule
    '22222222-2222-2222-2222-222222222222', -- auth_code（例）
    TRUE                    -- is_authenticated
),
-- 3件目: ハイブリッド数学セミナー（ハイブリッド）
(
    'ハイブリッド数学セミナー',
    '数学振興協会',
    '2024-03-14',            -- start_date
    '09:00:00',             -- start_time
    '2024-03-14',            -- end_date
    '18:00:00',             -- end_time
    'hybrid@mathforum.example.com', -- email（サンプル）
    '大阪府',               -- prefecture（会場が大阪府）
    'ハイブリッド',         -- event_type
    TRUE,                   -- is_online（オンラインでも参加可）
    TRUE,                   -- is_offline（オフライン会場もあり）
    'https://example.com/hybrid-seminar', -- official_url
    NULL,                   -- online_lecture_url（必要に応じて追加）
    '大阪府（オンラインでも参加可）', -- venue
    'オフラインとオンラインの両方で参加できる数学セミナーです。', -- description
    '["Business","Vacation"]',                   -- tags
    '[{"name":"Heidi","title":"CEO","organization":"GlobalCorp"}]',                   -- speakers
    '[{"time":"19:00","title":"Welcome Speech","speaker":"Heidi"}]',                   -- schedule
    '33333333-3333-3333-3333-333333333333', -- auth_code（例）
    TRUE                    -- is_authenticated
);
