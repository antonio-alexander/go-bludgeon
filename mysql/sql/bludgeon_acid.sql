USE bludgeon;

ALTER TABLE timers ADD FOREIGN KEY IF NOT EXISTS fk_employee_id (employee_id)
    REFERENCES employees(id) ON DELETE CASCADE;
