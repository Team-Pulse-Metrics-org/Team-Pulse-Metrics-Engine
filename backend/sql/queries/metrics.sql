-- Total commits for a user
SELECT COUNT(*)
FROM activities
WHERE user_id = $1
  AND activity_type = 'commit';