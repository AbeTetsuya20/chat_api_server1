INSERT INTO
    `administrators` (id, created_at, updated_at)
SELECT
    'admin',
    NOW(),
    NOW()
WHERE
    NOT EXISTS (
            SELECT
                1
            FROM
                `administrators`
            WHERE
                    `id` = 'admin'
        );
