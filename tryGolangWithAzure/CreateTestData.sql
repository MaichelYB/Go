-- CREATE SCHEMA TestSchema;
-- GO

-- CREATE TABLE TestSchema.Users
-- (
--     Id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
--     Name NVARCHAR(50),
--     Job NVARCHAR(50),
--     Date DATE
-- );
-- GO

-- INSERT INTO TestSchema.Users
--     (Name, Job, Date)
-- VALUES
--     (N'Jared', N'Programmer', getdate()),
--     (N'Nikita', N'Designer', getdate()),
--     (N'Tom', N'System Analyst', getdate());
-- GO

SELECT *
FROM TestSchema.Users;
GO