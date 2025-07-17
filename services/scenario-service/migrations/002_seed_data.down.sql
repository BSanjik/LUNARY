DELETE FROM steps WHERE scenario_id IN (1, 2);
DELETE FROM scenarios WHERE id IN (1, 2);

-- Сбрасываем последовательности
SELECT setval('scenarios_id_seq', 0);
SELECT setval('steps_id_seq', 0); 