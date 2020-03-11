package query

// GetAutoIncrementStatement statement for getting auto increment
const GetAutoIncrementStatement = `SELECT COLUMN_NAME
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_NAME = '%s'
AND TABLE_SCHEMA = '%s'
AND EXTRA like 'auto_increment'
`

// GetConstraintStatement statement for getting reference table
const GetConstraintStatement = `SELECT TABLE_NAME, COLUMN_NAME
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE
REFERENCED_TABLE_SCHEMA = '%s'
AND REFERENCED_COLUMN_NAME = '%s'
AND REFERENCED_TABLE_NAME = '%s'
`

// ShowTableStatement show tables statement
const ShowTableStatement = `show tables`
