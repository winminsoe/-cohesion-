## ~~cohesion~~

- Simple implementation of table dependency viewer.
- Better to view with the [jq](https://stedolan.github.io/jq/)

### Example
> DB_DSN=root:root@/db_name go run main.go | jq .
```
{
  "table": "flow",
  "dependent_tables": [
    {
      "table": "flow_history",
      "column": "fk_flow"
    },
    {
      "table": "flow_item",
      "column": "fk_flow"
    }
  ]
}
{
  "table": "flow_history",
  "dependent_tables": null
}
{
  "table": "flow_item",
  "dependent_tables": null
}
{
  "table": "flow_status",
  "dependent_tables": [
    {
      "table": "flow",
      "column": "fk_flow_status"
    },
    {
      "table": "flow",
      "column": "retry_status"
    },
    {
      "table": "flow_history",
      "column": "from"
    },
    {
      "table": "flow_history",
      "column": "to"
    }
  ]
}
{
  "table": "users",
  "dependent_tables": [
    {
      "table": "flow_history",
      "column": "fk_user"
    }
  ]
}
```