-- Seed data for Mini EVV Logger
-- Schedules with a mix of past (missed/completed), today (upcoming/in_progress), and future dates

-- Past schedules (yesterday - should be missed or completed)
INSERT INTO schedules (client_name, start_at, end_at, location, latitude, longitude)
VALUES
    ('Eleanor Thompson', CURRENT_DATE - INTERVAL '2 days' + TIME '08:00', CURRENT_DATE - INTERVAL '2 days' + TIME '10:00', '123 Maple Street, Springfield', 39.792, -89.658),
    ('Robert Chen', CURRENT_DATE - INTERVAL '2 days' + TIME '11:00', CURRENT_DATE - INTERVAL '2 days' + TIME '13:00', '456 Oak Avenue, Springfield', 39.801, -89.645),
    ('Margaret Williams', CURRENT_DATE - INTERVAL '1 day' + TIME '09:00', CURRENT_DATE - INTERVAL '1 day' + TIME '11:00', '789 Pine Road, Springfield', 39.789, -89.650),
    ('James Anderson', CURRENT_DATE - INTERVAL '1 day' + TIME '14:00', CURRENT_DATE - INTERVAL '1 day' + TIME '16:00', '321 Elm Boulevard, Springfield', 39.795, -89.660);

-- Today's schedules
INSERT INTO schedules (client_name, start_at, end_at, location, latitude, longitude)
VALUES
    ('Patricia Davis', CURRENT_DATE + TIME '08:00', CURRENT_DATE + TIME '10:00', '555 Cedar Lane, Springfield', 39.790, -89.655),
    ('William Johnson', CURRENT_DATE + TIME '11:00', CURRENT_DATE + TIME '13:00', '777 Birch Court, Springfield', 39.805, -89.640),
    ('Susan Martinez', CURRENT_DATE + TIME '14:00', CURRENT_DATE + TIME '16:00', '999 Walnut Drive, Springfield', 39.780, -89.670);

-- Future schedules
INSERT INTO schedules (client_name, start_at, end_at, location, latitude, longitude)
VALUES
    ('Dorothy Brown', CURRENT_DATE + INTERVAL '1 day' + TIME '09:00', CURRENT_DATE + INTERVAL '1 day' + TIME '11:00', '111 Ash Street, Springfield', 39.785, -89.665),
    ('Richard Wilson', CURRENT_DATE + INTERVAL '1 day' + TIME '13:00', CURRENT_DATE + INTERVAL '1 day' + TIME '15:00', '222 Poplar Poplar Way, Springfield', 39.798, -89.648),
    ('Karen Taylor', CURRENT_DATE + INTERVAL '2 days' + TIME '10:00', CURRENT_DATE + INTERVAL '2 days' + TIME '12:00', '333 Cypress Ave, Springfield', 39.810, -89.635);

-- Tasks for Eleanor Thompson (completed schedule)
INSERT INTO tasks (schedule_id, title, status) VALUES
    (1, 'Check vital signs (blood pressure, temperature)', 'completed'),
    (1, 'Administer morning medication', 'completed'),
    (1, 'Assist with breakfast preparation', 'completed'),
    (1, 'Light physical therapy exercises', 'completed');

-- Tasks for Robert Chen (missed schedule)
INSERT INTO tasks (schedule_id, title, status) VALUES
    (2, 'Wound care and bandage change', 'pending'),
    (2, 'Medication administration', 'pending'),
    (2, 'Mobility assistance', 'pending');

-- Tasks for Margaret Williams (completed schedule)
INSERT INTO tasks (schedule_id, title, status) VALUES
    (3, 'Assist with bathing and grooming', 'completed'),
    (3, 'Prepare and serve lunch', 'completed'),
    (3, 'Medication reminder and supervision', 'completed'),
    (3, 'Document health observations', 'completed');

-- Tasks for James Anderson (missed schedule)
INSERT INTO tasks (schedule_id, title, status) VALUES
    (4, 'Physical therapy session', 'pending'),
    (4, 'Meal preparation', 'pending'),
    (4, 'Medication administration', 'pending'),
    (4, 'Social engagement activities', 'pending');

-- Tasks for Patricia Davis (today - upcoming)
INSERT INTO tasks (schedule_id, title, status) VALUES
    (5, 'Morning medication administration', 'pending'),
    (5, 'Assist with personal hygiene', 'pending'),
    (5, 'Prepare breakfast and monitor intake', 'pending'),
    (5, 'Check and record vital signs', 'pending');

-- Tasks for William Johnson (today - upcoming)
INSERT INTO tasks (schedule_id, title, status) VALUES
    (6, 'Assist with bathing', 'pending'),
    (6, 'Give prescribed medication', 'pending'),
    (6, 'Light housekeeping duties', 'pending'),
    (6, 'Accompany to garden for fresh air', 'pending'),
    (6, 'Document care activities', 'pending');

-- Tasks for Susan Martinez (today - upcoming)
INSERT INTO tasks (schedule_id, title, status) VALUES
    (7, 'Wound dressing change', 'pending'),
    (7, 'Administer afternoon medication', 'pending'),
    (7, 'Range of motion exercises', 'pending'),
    (7, 'Prepare and serve snack', 'pending');

-- Tasks for Dorothy Brown (tomorrow)
INSERT INTO tasks (schedule_id, title, status) VALUES
    (8, 'Morning routine assistance', 'pending'),
    (8, 'Medication administration', 'pending'),
    (8, 'Prepare meals for the day', 'pending');

-- Tasks for Richard Wilson (tomorrow)
INSERT INTO tasks (schedule_id, title, status) VALUES
    (9, 'Physical therapy exercises', 'pending'),
    (9, 'Assist with lunch', 'pending'),
    (9, 'Companionship and conversation', 'pending');

-- Tasks for Karen Taylor (day after tomorrow)
INSERT INTO tasks (schedule_id, title, status) VALUES
    (10, 'Vital signs monitoring', 'pending'),
    (10, 'Medication management', 'pending'),
    (10, 'Assist with daily activities', 'pending'),
    (10, 'Prepare and organize weekly pill box', 'pending');
