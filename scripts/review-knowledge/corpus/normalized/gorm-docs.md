---
title: Advanced Query
layout: page
---

## <span id="smart_select">Smart Select Fields</span>

In GORM, you can efficiently select specific fields using the [`Select`](query.html) method. This is particularly useful when dealing with large models but requiring only a subset of fields, especially in API responses.

```go
type User struct {
 ID uint
 Name string
 Age int
 Gender string
 // hundreds of fields
}

type APIUser struct {
 ID uint
 Name string
}

// GORM will automatically select `id`, `name` fields when querying
db.Model(&User{}).Limit(10).Find(&APIUser{})
// SQL: SELECT `id`, `name` FROM `users` LIMIT 10
```

{% note warn %}
**NOTE** In `QueryFields` mode, all model fields are selected by their names.
{% endnote %}

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 QueryFields: true,
})

// Default behavior with QueryFields set to true
db.Find(&user)
// SQL: SELECT `users`.`name`, `users`.`age`, ... FROM `users`

// Using Session Mode with QueryFields
db.Session(&gorm.Session{QueryFields: true}).Find(&user)
// SQL: SELECT `users`.`name`, `users`.`age`, ... FROM `users`
```

## Locking

GORM supports different types of locks, for example:

```go
// Basic FOR UPDATE lock
db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
// SQL: SELECT * FROM `users` FOR UPDATE
```

The above statement will lock the selected rows for the duration of the transaction. This can be used in scenarios where you are preparing to update the rows and want to prevent other transactions from modifying them until your transaction is complete.

The `Strength` can be also set to `SHARE` which locks the rows in a way that allows other transactions to read the locked rows but not to update or delete them.
```go
db.Clauses(clause.Locking{
 Strength: "SHARE",
 Table: clause.Table{Name: clause.CurrentTable},
}).Find(&users)
// SQL: SELECT * FROM `users` FOR SHARE OF `users`
```

The `Table` option can be used to specify the table to lock. This is useful when you are joining multiple tables and want to lock only one of them.

Options can be provided like `NOWAIT` which tries to acquire a lock and fails immediately with an error if the lock is not available. It prevents the transaction from waiting for other transactions to release their locks.

```go
db.Clauses(clause.Locking{
 Strength: "UPDATE",
 Options: "NOWAIT",
}).Find(&users)
// SQL: SELECT * FROM `users` FOR UPDATE NOWAIT
```

Another option can be `SKIP LOCKED` which skips over any rows that are already locked by other transactions. This is useful in high concurrency situations where you want to process rows that are not currently locked by other transactions.

For more advanced locking strategies, refer to [Raw SQL and SQL Builder](sql_builder.html).

## SubQuery

Subqueries are a powerful feature in SQL, allowing nested queries. GORM can generate subqueries automatically when using a *gorm.DB object as a parameter.

```go
// Simple subquery
db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)
// SQL: SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");

// Nested subquery
subQuery := db.Select("AVG(age)").Where("name LIKE ?", "name%").Table("users")
db.Select("AVG(age) as avgage").Group("name").Having("AVG(age) > (?)", subQuery).Find(&results)
// SQL: SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%")
```

### <span id="from_subquery">From SubQuery</span>

GORM allows the use of subqueries in the FROM clause, enabling complex queries and data organization.

```go
// Using subquery in FROM clause
db.Table("(?) as u", db.Model(&User{}).Select("name", "age")).Where("age = ?", 18).Find(&User{})
// SQL: SELECT * FROM (SELECT `name`,`age` FROM `users`) as u WHERE `age` = 18

// Combining multiple subqueries in FROM clause
subQuery1 := db.Model(&User{}).Select("name")
subQuery2 := db.Model(&Pet{}).Select("name")
db.Table("(?) as u, (?) as p", subQuery1, subQuery2).Find(&User{})
// SQL: SELECT * FROM (SELECT `name` FROM `users`) as u, (SELECT `name` FROM `pets`) as p
```

## <span id="group_conditions">Group Conditions</span>

Group Conditions in GORM provide a more readable and maintainable way to write complex SQL queries involving multiple conditions.

```go
// Complex SQL query using Group Conditions
db.Where(
 db.Where("pizza = ?", "pepperoni").Where(db.Where("size = ?", "small").Or("size = ?", "medium")),
).Or(
 db.Where("pizza = ?", "hawaiian").Where("size = ?", "xlarge"),
).Find(&Pizza{})
// SQL: SELECT * FROM `pizzas` WHERE (pizza = "pepperoni" AND (size = "small" OR size = "medium")) OR (pizza = "hawaiian" AND size = "xlarge")
```

## IN with multiple columns

GORM supports the IN clause with multiple columns, allowing you to filter data based on multiple field values in a single query.

```go
// Using IN with multiple columns
db.Where("(name, age, role) IN ?", [][]interface{}{{"jinzhu", 18, "admin"}, {"jinzhu2", 19, "user"}}).Find(&users)
// SQL: SELECT * FROM users WHERE (name, age, role) IN (("jinzhu", 18, "admin"), ("jinzhu 2", 19, "user"));
```

## Named Argument

GORM enhances the readability and maintainability of SQL queries by supporting named arguments. This feature allows for clearer and more organized query construction, especially in complex queries with multiple parameters. Named arguments can be utilized using either [`sql.NamedArg`](https://tip.golang.org/pkg/database/sql/#NamedArg) or `map[string]interface{}{}`, providing flexibility in how you structure your queries.

```go
// Example using sql.NamedArg for named arguments
db.Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(&user)
// SQL: SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"

// Example using a map for named arguments
db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu"}).First(&user)
// SQL: SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu" ORDER BY `users`.`id` LIMIT 1
```

For more examples and details, see [Raw SQL and SQL Builder](sql_builder.html#named_argument)

## Find To Map

GORM provides flexibility in querying data by allowing results to be scanned into a `map[string]interface{}` or `[]map[string]interface{}`, which can be useful for dynamic data structures.

When using `Find To Map`, it's crucial to include `Model` or `Table` in your query to explicitly specify the table name. This ensures that GORM understands which table to query against.

```go
// Scanning the first result into a map with Model
result := map[string]interface{}{}
db.Model(&User{}).First(&result, "id = ?", 1)
// SQL: SELECT * FROM `users` WHERE id = 1 LIMIT 1

// Scanning multiple results into a slice of maps with Table
var results []map[string]interface{}
db.Table("users").Find(&results)
// SQL: SELECT * FROM `users`
```

## FirstOrInit

GORM's `FirstOrInit` method is utilized to fetch the first record that matches given conditions, or initialize a new instance if no matching record is found. This method is compatible with both struct and map conditions and allows additional flexibility with the `Attrs` and `Assign` methods.

```go
// If no User with the name "non_existing" is found, initialize a new User
var user User
db.FirstOrInit(&user, User{Name: "non_existing"})
// user -> User{Name: "non_existing"} if not found

// Retrieving a user named "jinzhu"
db.Where(User{Name: "jinzhu"}).FirstOrInit(&user)
// user -> User{ID: 111, Name: "Jinzhu", Age: 18} if found

// Using a map to specify the search condition
db.FirstOrInit(&user, map[string]interface{}{"name": "jinzhu"})
// user -> User{ID: 111, Name: "Jinzhu", Age: 18} if found
```

### Using `Attrs` for Initialization

When no record is found, you can use `Attrs` to initialize a struct with additional attributes. These attributes are included in the new struct but are not used in the SQL query.

```go
// If no User is found, initialize with given conditions and additional attributes
db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrInit(&user)
// SQL: SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
// user -> User{Name: "non_existing", Age: 20} if not found

// If a User named "Jinzhu" is found, `Attrs` are ignored
db.Where(User{Name: "Jinzhu"}).Attrs(User{Age: 20}).FirstOrInit(&user)
// SQL: SELECT * FROM USERS WHERE name = 'Jinzhu' ORDER BY id LIMIT 1;
// user -> User{ID: 111, Name: "Jinzhu", Age: 18} if found
```

### Using `Assign` for Attributes

The `Assign` method allows you to set attributes on the struct regardless of whether the record is found or not. These attributes are set on the struct but are not used to build the SQL query and the final data won't be saved into the database.

```go
// Initialize with given conditions and Assign attributes, regardless of record existence
db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrInit(&user)
// user -> User{Name: "non_existing", Age: 20} if not found

// If a User named "Jinzhu" is found, update the struct with Assign attributes
db.Where(User{Name: "Jinzhu"}).Assign(User{Age: 20}).FirstOrInit(&user)
// SQL: SELECT * FROM USERS WHERE name = 'Jinzhu' ORDER BY id LIMIT 1;
// user -> User{ID: 111, Name: "Jinzhu", Age: 20} if found
```

`FirstOrInit`, along with `Attrs` and `Assign`, provides a powerful and flexible way to ensure a record exists and is initialized or updated with specific attributes in a single step.

## FirstOrCreate

`FirstOrCreate` in GORM is used to fetch the first record that matches given conditions or create a new one if no matching record is found. This method is effective with both struct and map conditions. The `RowsAffected` property is useful to determine the number of records created or updated.

```go
// Create a new record if not found
result := db.FirstOrCreate(&user, User{Name: "non_existing"})
// SQL: INSERT INTO "users" (name) VALUES ("non_existing");
// user -> User{ID: 112, Name: "non_existing"}
// result.RowsAffected // => 1 (record created)

// If the user is found, no new record is created
result = db.Where(User{Name: "jinzhu"}).FirstOrCreate(&user)
// user -> User{ID: 111, Name: "jinzhu", Age: 18}
// result.RowsAffected // => 0 (no record created)
```

### Using `Attrs` with FirstOrCreate

`Attrs` can be used to specify additional attributes for the new record if it is not found. These attributes are used for creation but not in the initial search query.

```go
// Create a new record with additional attributes if not found
db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
// SQL: SELECT * FROM users WHERE name = 'non_existing';
// SQL: INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
// user -> User{ID: 112, Name: "non_existing", Age: 20}

// If the user is found, `Attrs` are ignored
db.Where(User{Name: "jinzhu"}).Attrs(User{Age: 20}).FirstOrCreate(&user)
// SQL: SELECT * FROM users WHERE name = 'jinzhu';
// user -> User{ID: 111, Name: "jinzhu", Age: 18}
```

### Using `Assign` with FirstOrCreate

The `Assign` method sets attributes on the record regardless of whether it is found or not, and these attributes are saved back to the database.

```go
// Initialize and save new record with `Assign` attributes if not found
db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrCreate(&user)
// SQL: SELECT * FROM users WHERE name = 'non_existing';
// SQL: INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
// user -> User{ID: 112, Name: "non_existing", Age: 20}

// Update found record with `Assign` attributes
db.Where(User{Name: "jinzhu"}).Assign(User{Age: 20}).FirstOrCreate(&user)
// SQL: SELECT * FROM users WHERE name = 'jinzhu';
// SQL: UPDATE users SET age=20 WHERE id = 111;
// user -> User{ID: 111, Name: "Jinzhu", Age: 20}
```

## Optimizer/Index Hints

GORM includes support for optimizer and index hints, allowing you to influence the query optimizer's execution plan. This can be particularly useful in optimizing query performance or when dealing with complex queries.

Optimizer hints are directives that suggest how a database's query optimizer should execute a query. GORM facilitates the use of optimizer hints through the gorm.io/hints package.

```go
import "gorm.io/hints"

// Using an optimizer hint to set a maximum execution time
db.Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find(&User{})
// SQL: SELECT * /*+ MAX_EXECUTION_TIME(10000) */ FROM `users`
```

### Index Hints

Index hints provide guidance to the database about which indexes to use. They can be beneficial if the query planner is not selecting the most efficient indexes for a query.

```go
import "gorm.io/hints"

// Suggesting the use of a specific index
db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
// SQL: SELECT * FROM `users` USE INDEX (`idx_user_name`)

// Forcing the use of certain indexes for a JOIN operation
db.Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find(&User{})
// SQL: SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)
```

These hints can significantly impact query performance and behavior, especially in large databases or complex data models. For more detailed information and additional examples, refer to [Optimizer Hints/Index/Comment](hints.html) in the GORM documentation.

## Iteration

GORM supports the iteration over query results using the `Rows` method. This feature is particularly useful when you need to process large datasets or perform operations on each record individually.

You can iterate through rows returned by a query, scanning each row into a struct. This method provides granular control over how each record is handled.

```go
rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Rows()
defer rows.Close()

for rows.Next() {
 var user User
 // ScanRows scans a row into a struct
 db.ScanRows(rows, &user)

 // Perform operations on each user
}
```

This approach is ideal for complex data processing that cannot be easily achieved with standard query methods.

## FindInBatches

`FindInBatches` allows querying and processing records in batches. This is especially useful for handling large datasets efficiently, reducing memory usage and improving performance.

With `FindInBatches`, GORM processes records in specified batch sizes. Inside the batch processing function, you can apply operations to each batch of records.

```go
// Processing records in batches of 100
result := db.Where("processed = ?", false).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
 for _, result := range results {
 // Operations on each record in the batch
 }

 // Save changes to the records in the current batch
 tx.Save(&results)

 // tx.RowsAffected provides the count of records in the current batch
 // The variable 'batch' indicates the current batch number

 // Returning an error will stop further batch processing
 return nil
})

// result.Error contains any errors encountered during batch processing
// result.RowsAffected provides the count of all processed records across batches
```

`FindInBatches` is an effective tool for processing large volumes of data in manageable chunks, optimizing resource usage and performance.

## Query Hooks

GORM offers the ability to use hooks, such as `AfterFind`, which are triggered during the lifecycle of a query. These hooks allow for custom logic to be executed at specific points, such as after a record has been retrieved from the database.

This hook is useful for post-query data manipulation or default value settings. For more detailed information and additional hook types, refer to [Hooks](hooks.html) in the GORM documentation.

```go
func (u *User) AfterFind(tx *gorm.DB) (err error) {
 // Custom logic after finding a user
 if u.Role == "" {
 u.Role = "user" // Set default role if not specified
 }
 return
}

// Usage of AfterFind hook happens automatically when a User is queried
```

## <span id="pluck">Pluck</span>

The `Pluck` method in GORM is used to query a single column from the database and scan the result into a slice. This method is ideal for when you need to retrieve specific fields from a model.

If you need to query more than one column, you can use `Select` with [Scan](query.html) or [Find](query.html) instead.

```go
// Retrieving ages of all users
var ages []int64
db.Model(&User{}).Pluck("age", &ages)

// Retrieving names of all users
var names []string
db.Model(&User{}).Pluck("name", &names)

// Retrieving names from a different table
db.Table("deleted_users").Pluck("name", &names)

// Using Distinct with Pluck
db.Model(&User{}).Distinct().Pluck("Name", &names)
// SQL: SELECT DISTINCT `name` FROM `users`

// Querying multiple columns
db.Select("name", "age").Scan(&users)
db.Select("name", "age").Find(&users)
```

## Scopes

`Scopes` in GORM are a powerful feature that allows you to define commonly-used query conditions as reusable methods. These scopes can be easily referenced in your queries, making your code more modular and readable.

### Defining Scopes

`Scopes` are defined as functions that modify and return a `gorm.DB` instance. You can define a variety of conditions as scopes based on your application's requirements.

```go
// Scope for filtering records where amount is greater than 1000
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
 return db.Where("amount > ?", 1000)
}

// Scope for orders paid with a credit card
func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
 return db.Where("pay_mode_sign = ?", "C")
}

// Scope for orders paid with cash on delivery (COD)
func PaidWithCod(db *gorm.DB) *gorm.DB {
 return db.Where("pay_mode_sign = ?", "COD")
}

// Scope for filtering orders by status
func OrderStatus(status []string) func(db *gorm.DB) *gorm.DB {
 return func(db *gorm.DB) *gorm.DB {
 return db.Where("status IN (?)", status)
 }
}
```

### Applying Scopes in Queries

You can apply one or more scopes to a query by using the `Scopes` method. This allows you to chain multiple conditions dynamically.

```go
// Applying scopes to find all credit card orders with an amount greater than 1000
db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)

// Applying scopes to find all COD orders with an amount greater than 1000
db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&orders)

// Applying scopes to find all orders with specific statuses and an amount greater than 1000
db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
```

`Scopes` are a clean and efficient way to encapsulate common query logic, enhancing the maintainability and readability of your code. For more detailed examples and usage, refer to [Scopes](scopes.html) in the GORM documentation.

## <span id="count">Count</span>

The `Count` method in GORM is used to retrieve the number of records that match a given query. It's a useful feature for understanding the size of a dataset, particularly in scenarios involving conditional queries or data analysis.

### Getting the Count of Matched Records

You can use `Count` to determine the number of records that meet specific criteria in your queries.

```go
var count int64

// Counting users with specific names
db.Model(&User{}).Where("name = ?", "jinzhu").Or("name = ?", "jinzhu 2").Count(&count)
// SQL: SELECT count(1) FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2'

// Counting users with a single name condition
db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count)
// SQL: SELECT count(1) FROM users WHERE name = 'jinzhu'

// Counting records in a different table
db.Table("deleted_users").Count(&count)
// SQL: SELECT count(1) FROM deleted_users
```

### Count with Distinct and Group

GORM also allows counting distinct values and grouping results.

```go
// Counting distinct names
db.Model(&User{}).Distinct("name").Count(&count)
// SQL: SELECT COUNT(DISTINCT(`name`)) FROM `users`

// Counting distinct values with a custom select
db.Table("deleted_users").Select("count(distinct(name))").Count(&count)
// SQL: SELECT count(distinct(name)) FROM deleted_users

// Counting grouped records
users := []User{
 {Name: "name1"},
 {Name: "name2"},
 {Name: "name3"},
 {Name: "name3"},
}

db.Model(&User{}).Group("name").Count(&count)
// Count after grouping by name
// count => 3
```

---
title: Associations
layout: page
---

## Auto Create/Update

GORM automates the saving of associations and their references when creating or updating records, using an upsert technique that primarily updates foreign key references for existing associations.

### Auto-Saving Associations on Create

When you create a new record, GORM will automatically save its associated data. This includes inserting data into related tables and managing foreign key references.

```go
user := User{
 Name: "jinzhu",
 BillingAddress: Address{Address1: "Billing Address - Address 1"},
 ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
 Emails: []Email{
 {Email: "jinzhu@example.com"},
 {Email: "jinzhu-2@example.com"},
 },
 Languages: []Language{
 {Name: "ZH"},
 {Name: "EN"},
 },
}

// Creating a user along with its associated addresses, emails, and languages
db.Create(&user)
// BEGIN TRANSACTION;
// INSERT INTO "addresses" (address1) VALUES ("Billing Address - Address 1"), ("Shipping Address - Address 1") ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "users" (name,billing_address_id,shipping_address_id) VALUES ("jinzhu", 1, 2);
// INSERT INTO "emails" (user_id,email) VALUES (111, "jinzhu@example.com"), (111, "jinzhu-2@example.com") ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "languages" ("name") VALUES ('ZH'), ('EN') ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "user_languages" ("user_id","language_id") VALUES (111, 1), (111, 2) ON DUPLICATE KEY DO NOTHING;
// COMMIT;

db.Save(&user)
```

### Updating Associations with `FullSaveAssociations`

For scenarios where a full update of the associated data is required (not just the foreign key references), the `FullSaveAssociations` mode should be used.

```go
// Update a user and fully update all its associations
db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)
// SQL: Fully updates addresses, users, emails tables, including existing associated records
```

Using `FullSaveAssociations` ensures that the entire state of the model, including all its associations, is reflected in the database, maintaining data integrity and consistency throughout the application.

## Skip Auto Create/Update

GORM provides flexibility to skip automatic saving of associations during create or update operations. This can be achieved using the `Select` or `Omit` methods, which allow you to specify exactly which fields or associations should be included or excluded in the operation.

### Using `Select` to Include Specific Fields

The `Select` method lets you specify which fields of the model should be saved. This means that only the selected fields will be included in the SQL operation.

```go
user := User{
 // User and associated data
}

// Only include the 'Name' field when creating the user
db.Select("Name").Create(&user)
// SQL: INSERT INTO "users" (name) VALUES ("jinzhu");
```

### Using `Omit` to Exclude Fields or Associations

Conversely, `Omit` allows you to exclude certain fields or associations when saving a model.

```go
// Skip creating the 'BillingAddress' when creating the user
db.Omit("BillingAddress").Create(&user)

// Skip all associations when creating the user
db.Omit(clause.Associations).Create(&user)
```

{% note warn %}
**NOTE:**
For many-to-many associations, GORM upserts the associations before creating join table references. To skip this upserting, use `Omit` with the association name followed by `.*`:

```go
// Skip upserting 'Languages' associations
db.Omit("Languages.*").Create(&user)
```

To skip creating both the association and its references:

```go
// Skip creating 'Languages' associations and their references
db.Omit("Languages").Create(&user)
```
{% endnote %}

Using `Select` and `Omit`, you can fine-tune how GORM handles the creation or updating of your models, giving you control over the auto-save behavior of associations.

## Select/Omit Association fields

In GORM, when creating or updating records, you can use the `Select` and `Omit` methods to specifically include or exclude certain fields of an associated model.

With `Select`, you can specify which fields of an associated model should be included when saving the primary model. This is particularly useful for selectively saving parts of an association.

Conversely, `Omit` lets you exclude certain fields of an associated model from being saved. This can be useful when you want to prevent specific parts of an association from being persisted.

```go
user := User{
 Name: "jinzhu",
 BillingAddress: Address{Address1: "Billing Address - Address 1", Address2: "addr2"},
 ShippingAddress: Address{Address1: "Shipping Address - Address 1", Address2: "addr2"},
}

// Create user and his BillingAddress, ShippingAddress, including only specified fields of BillingAddress
db.Select("BillingAddress.Address1", "BillingAddress.Address2").Create(&user)
// SQL: Creates user and BillingAddress with only 'Address1' and 'Address2' fields

// Create user and his BillingAddress, ShippingAddress, excluding specific fields of BillingAddress
db.Omit("BillingAddress.Address2", "BillingAddress.CreatedAt").Create(&user)
// SQL: Creates user and BillingAddress, omitting 'Address2' and 'CreatedAt' fields
```

## Delete Associations

GORM allows for the deletion of specific associated relationships (has one, has many, many2many) using the `Select` method when deleting a primary record. This feature is particularly useful for maintaining database integrity and ensuring related data is appropriately managed upon deletion.

You can specify which associations should be deleted along with the primary record by using `Select`.

```go
// Delete a user's account when deleting the user
db.Select("Account").Delete(&user)

// Delete a user's Orders and CreditCards associations when deleting the user
db.Select("Orders", "CreditCards").Delete(&user)

// Delete all of a user's has one, has many, and many2many associations
db.Select(clause.Associations).Delete(&user)

// Delete each user's account when deleting multiple users
db.Select("Account").Delete(&users)
```

{% note warn %}
**NOTE:**
It's important to note that associations will be deleted only if the primary key of the deleting record is not zero. GORM uses these primary keys as conditions to delete the selected associations.

```go
// This will not work as intended
db.Select("Account").Where("name = ?", "jinzhu").Delete(&User{})
// SQL: Deletes all users with name 'jinzhu', but their accounts won't be deleted

// Correct way to delete a user and their account
db.Select("Account").Where("name = ?", "jinzhu").Delete(&User{ID: 1})
// SQL: Deletes the user with name 'jinzhu' and ID '1', and the user's account

// Deleting a user with a specific ID and their account
db.Select("Account").Delete(&User{ID: 1})
// SQL: Deletes the user with ID '1', and the user's account
```
{% endnote %}

## Association Mode

Association Mode in GORM offers various helper methods to handle relationships between models, providing an efficient way to manage associated data.

To start Association Mode, specify the source model and the relationship's field name. The source model must contain a primary key, and the relationship's field name should match an existing association.

```go
var user User
db.Model(&user).Association("Languages")
// Check for errors
error := db.Model(&user).Association("Languages").Error
```

### Finding Associations

Retrieve associated records with or without additional conditions.

```go
// Simple find
db.Model(&user).Association("Languages").Find(&languages)

// Find with conditions
codes := []string{"zh-CN", "en-US", "ja-JP"}
db.Model(&user).Where("code IN ?", codes).Association("Languages").Find(&languages)
```

### Appending Associations

Add new associations for `many to many`, `has many`, or replace the current association for `has one`, `belongs to`.

```go
// Append new languages
db.Model(&user).Association("Languages").Append([]Language{languageZH, languageEN})

db.Model(&user).Association("Languages").Append(&Language{Name: "DE"})

db.Model(&user).Association("CreditCard").Append(&CreditCard{Number: "411111111111"})
```

### Replacing Associations

Replace current associations with new ones.

```go
// Replace existing languages
db.Model(&user).Association("Languages").Replace([]Language{languageZH, languageEN})

db.Model(&user).Association("Languages").Replace(Language{Name: "DE"}, languageEN)
```

### Deleting Associations

Remove the relationship between the source and arguments, only deleting the reference.

```go
// Delete specific languages
db.Model(&user).Association("Languages").Delete([]Language{languageZH, languageEN})

db.Model(&user).Association("Languages").Delete(languageZH, languageEN)
```

### Clearing Associations

Remove all references between the source and association.

```go
// Clear all languages
db.Model(&user).Association("Languages").Clear()
```

### Counting Associations

Get the count of current associations, with or without conditions.

```go
// Count all languages
db.Model(&user).Association("Languages").Count()

// Count with conditions
codes := []string{"zh-CN", "en-US", "ja-JP"}
db.Model(&user).Where("code IN ?", codes).Association("Languages").Count()
```

### Batch Data Handling

Association Mode allows you to handle relationships for multiple records in a batch. This includes finding, appending, replacing, deleting, and counting operations for associated data.

- **Finding Associations**: Retrieve associated data for a collection of records.

```go
db.Model(&users).Association("Role").Find(&roles)
```

- **Deleting Associations**: Remove specific associations across multiple records.

```go
db.Model(&users).Association("Team").Delete(&userA)
```

- **Counting Associations**: Get the count of associations for a batch of records.

```go
db.Model(&users).Association("Team").Count()
```

- **Appending/Replacing Associations**: Manage associations for multiple records. Note the need for matching argument lengths with the data.

```go
var users = []User{user1, user2, user3}

// Append different teams to different users in a batch
// Append userA to user1's team, userB to user2's team, and userA, userB, userC to user3's team
db.Model(&users).Association("Team").Append(&userA, &userB, &[]User{userA, userB, userC})

// Replace teams for multiple users in a batch
// Reset user1's team to userA, user2's team to userB, and user3's team to userA, userB, and userC
db.Model(&users).Association("Team").Replace(&userA, &userB, &[]User{userA, userB, userC})
```

## <span id="delete_association_record">Delete Association Record</span>

In GORM, the `Replace`, `Delete`, and `Clear` methods in Association Mode primarily affect the foreign key references, not the associated records themselves. Understanding and managing this behavior is crucial for data integrity.

- **Reference Update**: These methods update the association's foreign key to null, effectively removing the link between the source and associated models.
- **No Physical Record Deletion**: The actual associated records remain untouched in the database.

### Modifying Deletion Behavior with `Unscoped`

For scenarios requiring actual deletion of associated records, the `Unscoped` method alters this behavior.

- **Soft Delete**: Marks associated records as deleted (sets `deleted_at` field) without removing them from the database.

```go
db.Model(&user).Association("Languages").Unscoped().Clear()
```

- **Permanent Delete**: Physically deletes the association records from the database.

```go
// db.Unscoped().Model(&user)
db.Unscoped().Model(&user).Association("Languages").Unscoped().Clear()
```

## <span id="tags">Association Tags</span>

Association tags in GORM are used to specify how associations between models are handled. These tags define the relationship's details, such as foreign keys, references, and constraints. Understanding these tags is essential for setting up and managing relationships effectively.

| Tag | Description |
|------------------|-----------------------------------------------------------------------------------------------------|
| `foreignKey` | Specifies the column name of the current model used as a foreign key in the join table. |
| `references` | Indicates the column name in the reference table that the foreign key of the join table maps to. |
| `polymorphic` | Defines the polymorphic type, typically the model name. |
| `polymorphicValue` | Sets the polymorphic value, usually the table name, if not specified otherwise. |
| `many2many` | Names the join table used in a many-to-many relationship. |
| `joinForeignKey` | Identifies the foreign key column in the join table that maps back to the current model's table. |
| `joinReferences` | Points to the foreign key column in the join table that links to the reference model's table. |
| `constraint` | Specifies relational constraints like `OnUpdate`, `OnDelete` for the association. |

---
title: Belongs To
layout: page
---

## Belongs To

A `belongs to` association sets up a one-to-one connection with another model, such that each instance of the declaring model "belongs to" one instance of the other model.

For example, if your application includes users and companies, and each user can be assigned to exactly one company, the following types represent that relationship. Notice here that, on the `User` object, there is both a `CompanyID` as well as a `Company`. By default, the `CompanyID` is implicitly used to create a foreign key relationship between the `User` and `Company` tables, and thus must be included in the `User` struct in order to fill the `Company` inner struct.

```go
// `User` belongs to `Company`, `CompanyID` is the foreign key
type User struct {
 gorm.Model
 Name string
 CompanyID int
 Company Company
}

type Company struct {
 ID int
 Name string
}
```

Refer to [Eager Loading](belongs_to.html#Eager-Loading) for details on populating the inner struct.

## Override Foreign Key

To define a belongs to relationship, the foreign key must exist, the default foreign key uses the owner's type name plus its primary field name.

For the above example, to define the `User` model that belongs to `Company`, the foreign key should be `CompanyID` by convention

GORM provides a way to customize the foreign key, for example:

```go
type User struct {
 gorm.Model
 Name string
 CompanyRefer int
 Company Company `gorm:"foreignKey:CompanyRefer"`
 // use CompanyRefer as foreign key
}

type Company struct {
 ID int
 Name string
}
```

## Override References

For a belongs to relationship, GORM usually uses the owner's primary field as the foreign key's value, for the above example, it is `Company`'s field `ID`.

When you assign a user to a company, GORM will save the company's `ID` into the user's `CompanyID` field.

You are able to change it with tag `references`, e.g:

```go
type User struct {
 gorm.Model
 Name string
 CompanyID string
 Company Company `gorm:"references:Code"` // use Code as references
}

type Company struct {
 ID int
 Code string
 Name string
}
```

{% note warn %}
**NOTE** GORM usually guess the relationship as `has one` if override foreign key name already exists in owner's type, we need to specify `references` in the `belongs to` relationship.
{% endnote %}

```go
type User struct {
 gorm.Model
 Name string
 CompanyID int
 Company Company `gorm:"references:CompanyID"` // use Company.CompanyID as references
}

type Company struct {
 CompanyID int
 Code string
 Name string
}
```

## CRUD with Belongs To

Please checkout [Association Mode](associations.html#Association-Mode) for working with belongs to relations

## Eager Loading

GORM allows eager loading belongs to associations with `Preload` or `Joins`, refer [Preloading (Eager loading)](preload.html) for details

## FOREIGN KEY Constraints

You can setup `OnUpdate`, `OnDelete` constraints with tag `constraint`, it will be created when migrating with GORM, for example:

```go
type User struct {
 gorm.Model
 Name string
 CompanyID int
 Company Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Company struct {
 ID int
 Name string
}
```

---
title: Change Log
layout: page
---

## v2.0 - 2020.08

GORM 2.0 is a rewrite from scratch, it introduces some incompatible-API change and many improvements

* Performance Improvements
* Modularity
* Context, Batch Insert, Prepared Statement Mode, DryRun Mode, Join Preload, Find To Map, Create From Map, FindInBatches supports
* Nested Transaction/SavePoint/RollbackTo SavePoint supports
* Named Argument, Group Conditions, Upsert, Locking, Optimizer/Index/Comment Hints supports, SubQuery improvements
* Full self-reference relationships supports, Join Table improvements, Association Mode for batch data
* Multiple fields support for tracking create/update time, which adds support for UNIX (milli/nano) seconds
* Field permissions support: read-only, write-only, create-only, update-only, ignored
* New plugin system: multiple databases, read/write splitting support with plugin Database Resolver, prometheus integrations...
* New Hooks API: unified interface with plugins
* New Migrator: allows to create database foreign keys for relationships, constraints/checker support, enhanced index support
* New Logger: context support, improved extensibility
* Unified Naming strategy: table name, field name, join table name, foreign key, checker, index name rules
* Better customized data type support (e.g: JSON)

[GORM 2.0 Release Note](v2_release_note.html)

## v1.0 - 2016.04

[GORM V1 Docs](https://v1.gorm.io)

Breaking Changes:

* `gorm.Open` returns `*gorm.DB` instead of `gorm.DB`
* Updating will only update changed fields
* Soft Delete's will only check `deleted_at IS NULL`
* New ToDBName logic
 Common initialisms from [golint](https://github.com/golang/lint/blob/master/lint.go#L702) like `HTTP`, `URI` was converted to lowercase, so `HTTP`'s db name is `http`, but not `h_t_t_p`, but for some other initialisms not in the list, like `SKU`, it's db name was `s_k_u`, this change fixed it to `sku`
* Error `RecordNotFound` has been renamed to `ErrRecordNotFound`
* `mssql` dialect has been renamed to `github.com/jinzhu/gorm/dialects/mssql`
* `Hstore` has been moved to package `github.com/jinzhu/gorm/dialects/postgres`

---
title: Composite Primary Key
layout: page
---

Set multiple fields as primary key creates composite primary key, for example:

```go
type Product struct {
 ID string `gorm:"primaryKey"`
 LanguageCode string `gorm:"primaryKey"`
 Code string
 Name string
}
```

**Note** integer `PrioritizedPrimaryField` enables `AutoIncrement` by default, to disable it, you need to turn off `autoIncrement` for the int fields:

```go
type Product struct {
 CategoryID uint64 `gorm:"primaryKey;autoIncrement:false"`
 TypeID uint64 `gorm:"primaryKey;autoIncrement:false"`
}
```

---
title: Connecting to a Database
layout: page
---

GORM officially supports the databases MySQL, PostgreSQL, GaussDB, SQLite, SQL Server TiDB, and Oracle Database

## MySQL

```go
import (
 "gorm.io/driver/mysql"
 "gorm.io/gorm"
)

func main() {
 // refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
 dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
 db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

{% note warn %}
**NOTE:**
To handle `time.Time` correctly, you need to include `parseTime` as a parameter. ([more parameters](https://github.com/go-sql-driver/mysql#parameters))
To fully support UTF-8 encoding, you need to change `charset=utf8` to `charset=utf8mb4`. See [this article](https://mathiasbynens.be/notes/mysql-utf8mb4) for a detailed explanation
{% endnote %}

MySQL Driver provides a [few advanced configurations](https://github.com/go-gorm/mysql) which can be used during initialization, for example:

```go
db, err := gorm.Open(mysql.New(mysql.Config{
 DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name
 DefaultStringSize: 256, // default size for string fields
 DisableDatetimePrecision: true, // disable datetime precision, which not supported before MySQL 5.6
 DontSupportRenameIndex: true, // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
 DontSupportRenameColumn: true, // `change` when rename column, rename column not supported before MySQL 8, MariaDB
 SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
}), &gorm.Config{})
```

### Customize Driver

GORM allows to customize the MySQL driver with the `DriverName` option, for example:

```go
import (
 _ "example.com/my_mysql_driver"
 "gorm.io/driver/mysql"
 "gorm.io/gorm"
)

db, err := gorm.Open(mysql.New(mysql.Config{
 DriverName: "my_mysql_driver",
 DSN: "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
}), &gorm.Config{})
```

### Existing database connection

GORM allows to initialize `*gorm.DB` with an existing database connection

```go
import (
 "database/sql"
 "gorm.io/driver/mysql"
 "gorm.io/gorm"
)

sqlDB, err := sql.Open("mysql", "mydb_dsn")
gormDB, err := gorm.Open(mysql.New(mysql.Config{
 Conn: sqlDB,
}), &gorm.Config{})
```

## PostgreSQL

```go
import (
 "gorm.io/driver/postgres"
 "gorm.io/gorm"
)

dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

We are using [pgx](https://github.com/jackc/pgx) as postgres's database/sql driver, it enables prepared statement cache by default, to disable it:

```go
// https://github.com/go-gorm/postgres
db, err := gorm.Open(postgres.New(postgres.Config{
 DSN: "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
 PreferSimpleProtocol: true, // disables implicit prepared statement usage
}), &gorm.Config{})
```

### Customize Driver

GORM allows to customize the PostgreSQL driver with the `DriverName` option, for example:

```go
import (
 _ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
 "gorm.io/gorm"
)

db, err := gorm.Open(postgres.New(postgres.Config{
 DriverName: "cloudsqlpostgres",
 DSN: "host=project:region:instance user=postgres dbname=postgres password=password sslmode=disable",
})
```

### Existing database connection

GORM allows to initialize `*gorm.DB` with an existing database connection

```go
import (
 "database/sql"
 "gorm.io/driver/postgres"
 "gorm.io/gorm"
)

sqlDB, err := sql.Open("pgx", "mydb_dsn")
gormDB, err := gorm.Open(postgres.New(postgres.Config{
 Conn: sqlDB,
}), &gorm.Config{})
```

## GaussDB

```go
import (
 "gorm.io/driver/gaussdb"
 "gorm.io/gorm"
)

dsn := "host=localhost user=gorm password=gorm dbname=gorm port=8000 sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(gaussdb.Open(dsn), &gorm.Config{})
```

We are using [gaussdb-go](https://github.com/HuaweiCloudDeveloper/gaussdb-go) as gaussdb's database/sql driver, it enables prepared statement cache by default, to disable it:

```go
// https://github.com/go-gorm/gaussdb
db, err := gorm.Open(gaussdb.New(gaussdb.Config{
 DSN: "user=gorm password=gorm dbname=gorm port=8000 sslmode=disable TimeZone=Asia/Shanghai",
 PreferSimpleProtocol: true, // disables implicit prepared statement usage
}), &gorm.Config{})
```

### Customize Driver

GORM allows to customize the GaussDB driver with the `DriverName` option, for example:

```go
import (
 _ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/gaussdb"
 "gorm.io/gorm"
)

db, err := gorm.Open(gaussdb.New(gaussdb.Config{
 DriverName: "cloudsqlgaussdb",
 DSN: "host=project:region:instance user=gaussdb dbname=gaussdb password=password sslmode=disable",
})
```

### Existing database connection

GORM allows to initialize `*gorm.DB` with an existing database connection

```go
import (
 "database/sql"
 "gorm.io/driver/gaussdb"
 "gorm.io/gorm"
)

sqlDB, err := sql.Open("gaussdbgo", "mydb_dsn")
gormDB, err := gorm.Open(gaussdb.New(gaussdb.Config{
 Conn: sqlDB,
}), &gorm.Config{})
```

## Oracle Database

The GORM Driver for Oracle provides support for Oracle Database, enabling full compatibility with GORM's ORM capabilities. It is built on top of the [Go Driver for Oracle (Godror)](https://github.com/godror/godror) and supports key features such as auto migrations, associations, transactions, and advanced querying.

### Prerequisite: Install Instant Client

To use ODPI-C with Godror, you’ll need to install the Oracle Instant Client on your system. Follow the steps on [this page](https://odpi-c.readthedocs.io/en/latest/user_guide/installation.html) to complete the installation.

After that, you can connect to the database using the `dataSourceName`, which specifies connection parameters (such as username and password) using a logfmt-encoded parameter list.

The way you specify the Instant Client directory differs by platform:

- macOS and Windows: You can set the `libDir` parameter in the dataSourceName.
- Linux: The libraries must be in the system library search path before your Go process starts, preferably configured with "ldconfig". The libDir parameter does not work on Linux.

#### Example (macOS/Windows)

``` go
dataSourceName := `user="scott" password="tiger"
 connectString="dbhost:1521/orclpdb1"
 libDir="/Path/to/your/instantclient_23_26"`
```

#### Example (Linux)

``` go
dataSourceName := `user="scott" password="tiger"
 connectString="dbhost:1521/orclpdb1"`
```

### Getting Started

```go
import (
 "github.com/oracle-samples/gorm-oracle/oracle"
 "gorm.io/gorm"
)

dataSourceName := `user="scott" password="tiger"
 connectString="dbhost:1521/orclpdb1"`
db, err := gorm.Open(oracle.Open(dataSourceName), &gorm.Config{})
```

## SQLite

```go
import (
 "gorm.io/driver/sqlite" // Sqlite driver based on CGO
 // "github.com/glebarez/sqlite" // Pure-Go SQLite driver, checkout https://github.com/glebarez/sqlite for details
 // "github.com/libtnb/sqlite" // Pure-Go SQLite driver, checkout https://github.com/libtnb/sqlite for details
 "gorm.io/gorm"
)

// github.com/mattn/go-sqlite3
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
```

{% note warn %}
**NOTE:** You can also use `file::memory:?cache=shared` instead of a path to a file. This will tell SQLite to use a temporary database in system memory. (See [SQLite docs](https://www.sqlite.org/inmemorydb.html) for this)
{% endnote %}

## SQL Server

```go
import (
 "gorm.io/driver/sqlserver"
 "gorm.io/gorm"
)

// github.com/denisenkom/go-mssqldb
dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
```

## TiDB

TiDB is compatible with MySQL protocol. You can follow the [MySQL](#mysql) part to create a connection to TiDB.

There are some points noteworthy for TiDB:

- You can use `gorm:"primaryKey;default:auto_random()"` tag to use [`AUTO_RANDOM`](https://docs.pingcap.com/tidb/stable/auto-random) feature for TiDB.
- TiDB supported [`SAVEPOINT`](https://docs.pingcap.com/tidb/stable/sql-statement-savepoint) from `v6.2.0`, please notice the version of TiDB when you use this feature.
- TiDB supported [`FOREIGN KEY`](https://docs.pingcap.com/tidb/dev/foreign-key) from `v6.6.0`, please notice the version of TiDB when you use this feature.

```go
import (
 "fmt"
 "gorm.io/driver/mysql"
 "gorm.io/gorm"
)

type Product struct {
 ID uint `gorm:"primaryKey;default:auto_random()"`
 Code string
 Price uint
}

func main() {
 db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:4000)/test"), &gorm.Config{})
 if err != nil {
 panic("failed to connect database")
 }

 db.AutoMigrate(&Product{})

 insertProduct := &Product{Code: "D42", Price: 100}

 db.Create(insertProduct)
 fmt.Printf("insert ID: %d, Code: %s, Price: %d\n",
 insertProduct.ID, insertProduct.Code, insertProduct.Price)

 readProduct := &Product{}
 db.First(&readProduct, "code = ?", "D42") // find product with code D42

 fmt.Printf("read ID: %d, Code: %s, Price: %d\n",
 readProduct.ID, readProduct.Code, readProduct.Price)
}
```

## Clickhouse

https://github.com/go-gorm/clickhouse

```go
import (
 "gorm.io/driver/clickhouse"
 "gorm.io/gorm"
)

func main() {
 dsn := "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
 db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})

 // Auto Migrate
 db.AutoMigrate(&User{})
 // Set table options
 db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&User{})

 // Insert
 db.Create(&user)

 // Select
 db.Find(&user, "id = ?", 10)

 // Batch Insert
 var users = []User{user1, user2, user3}
 db.Create(&users)
 // ...
}
```

## Connection Pool

GORM using [database/sql](https://pkg.go.dev/database/sql) to maintain connection pool

```go
sqlDB, err := db.DB()

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns sets the maximum number of open connections to the database.
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
sqlDB.SetConnMaxLifetime(time.Hour)
```

Refer [Generic Interface](generic_interface.html) for details

## Unsupported Databases

Some databases may be compatible with the `mysql` or `postgres` dialect, in which case you could just use the dialect for those databases.

For others, [you are encouraged to make a driver, pull request welcome!](write_driver.html)

---
title: Constraints
layout: page
---

GORM allows create database constraints with tag, constraints will be created when [AutoMigrate or CreateTable with GORM](migration.html)

## CHECK Constraint

Create CHECK constraints with tag `check`

```go
type UserIndex struct {
	Name string `gorm:"check:name_checker,name <> 'jinzhu'"`
	Name2 string `gorm:"check:name <> 'jinzhu'"`
	Name3 string `gorm:"check:,name <> 'jinzhu'"`
}
```

## Index Constraint

Checkout [Database Indexes](indexes.html)

## Foreign Key Constraint

GORM will creates foreign keys constraints for associations, you can disable this feature during initialization:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 DisableForeignKeyConstraintWhenMigrating: true,
})
```

GORM allows you setup FOREIGN KEY constraints's `OnDelete`, `OnUpdate` option with tag `constraint`, for example:

```go
type User struct {
 gorm.Model
 CompanyID int
 Company Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
 CreditCard CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreditCard struct {
 gorm.Model
 Number string
 UserID uint
}

type Company struct {
 ID int
 Name string
}
```

---
title: Context
layout: page
---

GORM's context support is a powerful feature that enhances the flexibility and control of database operations in Go applications. It allows for context management across different operational modes, timeout settings, and even integration into hooks/callbacks and middlewares. Let's delve into these various aspects:

## Single Session Mode

Single session mode is appropriate for executing individual operations. It ensures that the specific operation is executed within the context's scope, allowing for better control and monitoring.

### Generics API

With the Generics API, context is passed directly as the first parameter to the operation methods:

```go
users, err := gorm.G[User](db).Find(ctx)
```

### Traditional API

With the Traditional API, context is passed using the `WithContext` method:

```go
db.WithContext(ctx).Find(&users)
```

## Continuous Session Mode

Continuous session mode is ideal for performing a series of related operations. It maintains the context across these operations, which is particularly useful in scenarios like transactions.

```go
tx := db.WithContext(ctx)
tx.First(&user, 1)
tx.Model(&user).Update("role", "admin")
```

## Context Timeout

Setting a timeout on the context can control the duration of long-running queries. This is crucial for maintaining performance and avoiding resource lock-ups in database interactions.

### Generics API

With the Generics API, you pass the timeout context directly to the operation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

users, err := gorm.G[User](db).Find(ctx)
```

### Traditional API

With the Traditional API, you pass the timeout context to `WithContext`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

db.WithContext(ctx).Find(&users)
```

## Context in Hooks/Callbacks

The context can also be accessed within GORM's hooks/callbacks. This enables contextual information to be used during these lifecycle events. The context is accessible through the `Statement.Context` field:

```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
 ctx := tx.Statement.Context
 // ... use context
 return
}
```

## Integration with Chi Middleware

GORM's context support extends to web server middlewares, such as those in the Chi router. This allows setting a context with a timeout for all database operations within the scope of a web request.

```go
func SetDBMiddleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 timeoutContext, _ := context.WithTimeout(context.Background(), time.Second)
 ctx := context.WithValue(r.Context(), "DB", db.WithContext(timeoutContext))
 next.ServeHTTP(w, r.WithContext(ctx))
 })
}

// Router setup
r := chi.NewRouter()
r.Use(SetDBMiddleware)

// Route handlers
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
 db, ok := r.Context().Value("DB").(*gorm.DB)
 // ... db operations
})

r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
 db, ok := r.Context().Value("DB").(*gorm.DB)
 // ... db operations
})
```

**Note**: Setting the `Context` with `WithContext` is goroutine-safe. This ensures that database operations are safely managed across multiple goroutines. For more details, refer to the [Session documentation](session.html) in GORM.

## Logger Integration

GORM's logger also accepts `Context`, which can be used for log tracking and integrating with existing logging infrastructures.

Refer to [Logger documentation](logger.html) for more details.

---
title: Conventions
layout: page
---

## `ID` as Primary Key

GORM uses the field with the name `ID` as the table's primary key by default.

```go
type User struct {
  ID   string // field named `ID` will be used as a primary field by default
  Name string
}
```

You can set other fields as primary key with tag `primaryKey`

```go
// Set field `UUID` as primary field
type Animal struct {
  ID     int64
  UUID   string `gorm:"primaryKey"`
  Name   string
  Age    int64
}
```

Also check out [Composite Primary Key](composite_primary_key.html)

## Pluralized Table Name

GORM pluralizes struct name to `snake_cases` as table name, for struct `User`, its table name is `users` by convention

### TableName

You can change the default table name by implementing the `Tabler` interface, for example:

```go
type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `profiles`
func (User) TableName() string {
  return "profiles"
}
```

{% note warn %}
**NOTE** `TableName` doesn't allow dynamic name, its result will be cached for future, to use dynamic name, you can use `Scopes`, for example:
{% endnote %}

```go
func UserTable(user User) func (tx *gorm.DB) *gorm.DB {
 return func (tx *gorm.DB) *gorm.DB {
 if user.Admin {
 return tx.Table("admin_users")
 }

 return tx.Table("users")
 }
}

db.Scopes(UserTable(user)).Create(&user)
```

### Temporarily specify a name

Temporarily specify table name with `Table` method, for example:

```go
// Create table `deleted_users` with struct User's fields
db.Table("deleted_users").AutoMigrate(&User{})

// Query data from another table
var deletedUsers []User
db.Table("deleted_users").Find(&deletedUsers)
// SELECT * FROM deleted_users;

db.Table("deleted_users").Where("name = ?", "jinzhu").Delete(&User{})
// DELETE FROM deleted_users WHERE name = 'jinzhu';
```

Check out [From SubQuery](advanced_query.html#from_subquery) for how to use SubQuery in FROM clause

### <span id="naming_strategy">NamingStrategy</span>

GORM allows users to change the default naming conventions by overriding the default `NamingStrategy`, which is used to build `TableName`, `ColumnName`, `JoinTableName`, `RelationshipFKName`, `CheckerName`, `IndexName`, Check out [GORM Config](gorm_config.html#naming_strategy) for details

## Column Name

Column db name uses the field's name's `snake_case` by convention.

```go
type User struct {
  ID        uint      // column name is `id`
  Name      string    // column name is `name`
  Birthday  time.Time // column name is `birthday`
  CreatedAt time.Time // column name is `created_at`
}
```

You can override the column name with tag `column` or use [`NamingStrategy`](#naming_strategy)

```go
type Animal struct {
  AnimalID int64     `gorm:"column:beast_id"`         // set name to `beast_id`
  Birthday time.Time `gorm:"column:day_of_the_beast"` // set name to `day_of_the_beast`
  Age      int64     `gorm:"column:age_of_the_beast"` // set name to `age_of_the_beast`
}
```

## Timestamp Tracking

### CreatedAt

For models having `CreatedAt` field, the field will be set to the current time when the record is first created if its value is zero

```go
db.Create(&user) // set `CreatedAt` to current time

user2 := User{Name: "jinzhu", CreatedAt: time.Now()}
db.Create(&user2) // user2's `CreatedAt` won't be changed

// To change its value, you could use `Update`
db.Model(&user).Update("CreatedAt", time.Now())
```

You can disable the timestamp tracking by setting `autoCreateTime` tag to `false`, for example:

```go
type User struct {
 CreatedAt time.Time `gorm:"autoCreateTime:false"`
}
```

### UpdatedAt

For models having `UpdatedAt` field, the field will be set to the current time when the record is updated or created if its value is zero

```go
db.Save(&user) // set `UpdatedAt` to current time

db.Model(&user).Update("name", "jinzhu") // will set `UpdatedAt` to current time

db.Model(&user).UpdateColumn("name", "jinzhu") // `UpdatedAt` won't be changed

user2 := User{Name: "jinzhu", UpdatedAt: time.Now()}
db.Create(&user2) // user2's `UpdatedAt` won't be changed when creating

user3 := User{Name: "jinzhu", UpdatedAt: time.Now()}
db.Save(&user3) // user3's `UpdatedAt` will change to current time when updating
```

You can disable the timestamp tracking by setting `autoUpdateTime` tag to `false`, for example:

```go
type User struct {
 UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}
```

{% note %}
**NOTE** GORM supports having multiple time tracking fields and track with UNIX (nano/milli) seconds, checkout [Models](models.html#time_tracking) for more details
{% endnote %}

---
title: Create
layout: page
---

## Create Record

### Generics API

```go
user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

// Create a single record
ctx := context.Background()
err := gorm.G[User](db).Create(ctx, &user) // pass pointer of data to Create

// Create with result
result := gorm.WithResult()
err := gorm.G[User](db, result).Create(ctx, &user)
user.ID // returns inserted data's primary key
result.Error // returns error
result.RowsAffected // returns inserted records count
```

### Traditional API

```go
user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

result := db.Create(&user) // pass pointer of data to Create

user.ID // returns inserted data's primary key
result.Error // returns error
result.RowsAffected // returns inserted records count
```

We can also create multiple records with `Create()`:
```go
users := []*User{
	{Name: "Jinzhu", Age: 18, Birthday: time.Now()},
	{Name: "Jackson", Age: 19, Birthday: time.Now()},
}

result := db.Create(users) // pass a slice to insert multiple row

result.Error // returns error
result.RowsAffected // returns inserted records count
```

{% note warn %}
**NOTE** You cannot pass a struct to 'create', so you should pass a pointer to the data.
{% endnote %}

## Create Record With Selected Fields

Create a record and assign a value to the fields specified.

```go
db.Select("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`name`,`age`,`created_at`) VALUES ("jinzhu", 18, "2020-07-04 11:05:21.775")
```

Create a record and ignore the values for fields passed to omit.

```go
db.Omit("Name", "Age", "CreatedAt").Create(&user)
// INSERT INTO `users` (`birthday`,`updated_at`) VALUES ("2020-01-01 00:00:00.000", "2020-07-04 11:05:21.775")
```

## <span id="batch_insert">Batch Insert</span>

To efficiently insert large number of records, pass a slice to the `Create` method. GORM will generate a single SQL statement to insert all the data and backfill primary key values, hook methods will be invoked too. It will begin a **transaction** when records can be split into multiple batches.

```go
var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
db.Create(&users)

for _, user := range users {
 user.ID // 1,2,3
}
```

You can specify batch size when creating with `CreateInBatches`, e.g:

```go
var users = []User{{Name: "jinzhu_1"}, ...., {Name: "jinzhu_10000"}}

// batch size 100
db.CreateInBatches(users, 100)
```

Batch Insert is also supported when using [Upsert](#upsert) and [Create With Associations](#create_with_associations)

{% note warn %}
**NOTE** initialize GORM with `CreateBatchSize` option, all `INSERT` will respect this option when creating record & associations
{% endnote %}

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 CreateBatchSize: 1000,
})

db := db.Session(&gorm.Session{CreateBatchSize: 1000})

users = [5000]User{{Name: "jinzhu", Pets: []Pet{pet1, pet2, pet3}}...}

db.Create(&users)
// INSERT INTO users xxx (5 batches)
// INSERT INTO pets xxx (15 batches)
```

## Create Hooks

GORM allows user defined hooks to be implemented for `BeforeSave`, `BeforeCreate`, `AfterSave`, `AfterCreate`. These hook method will be called when creating a record, refer [Hooks](hooks.html) for details on the lifecycle

```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
 u.UUID = uuid.New()

	if u.Role == "admin" {
 return errors.New("invalid role")
	}
	return
}
```

If you want to skip `Hooks` methods, you can use the `SkipHooks` session mode, for example:

```go
DB.Session(&gorm.Session{SkipHooks: true}).Create(&user)

DB.Session(&gorm.Session{SkipHooks: true}).Create(&users)

DB.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(users, 100)
```

## Create From Map

GORM supports create from `map[string]interface{}` and `[]map[string]interface{}{}`, e.g:

```go
db.Model(&User{}).Create(map[string]interface{}{
 "Name": "jinzhu", "Age": 18,
})

// batch insert from `[]map[string]interface{}{}`
db.Model(&User{}).Create([]map[string]interface{}{
 {"Name": "jinzhu_1", "Age": 18},
 {"Name": "jinzhu_2", "Age": 20},
})
```

{% note warn %}
**NOTE** When creating from map, hooks won't be invoked, associations won't be saved and primary key values won't be back filled
{% endnote %}

## <span id="create_from_sql_expr">Create From SQL Expression/Context Valuer</span>

GORM allows insert data with SQL expression, there are two ways to achieve this goal, create from `map[string]interface{}` or [Customized Data Types](data_types.html#gorm_valuer_interface), for example:

```go
// Create from map
db.Model(User{}).Create(map[string]interface{}{
 "Name": "jinzhu",
 "Location": clause.Expr{SQL: "ST_PointFromText(?)", Vars: []interface{}{"POINT(100 100)"}},
})
// INSERT INTO `users` (`name`,`location`) VALUES ("jinzhu",ST_PointFromText("POINT(100 100)"));

// Create from customized data type
type Location struct {
	X, Y int
}

// Scan implements the sql.Scanner interface
func (loc *Location) Scan(v interface{}) error {
 // Scan a value into struct from database driver
}

func (loc Location) GormDataType() string {
 return "geometry"
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
 return clause.Expr{
 SQL: "ST_PointFromText(?)",
 Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y)},
 }
}

type User struct {
 Name string
 Location Location
}

db.Create(&User{
 Name: "jinzhu",
 Location: Location{X: 100, Y: 100},
})
// INSERT INTO `users` (`name`,`location`) VALUES ("jinzhu",ST_PointFromText("POINT(100 100)"))
```

## Advanced

### <span id="create_with_associations">Create With Associations</span>

When creating some data with associations, if its associations value is not zero-value, those associations will be upserted, and its `Hooks` methods will be invoked.

```go
type CreditCard struct {
 gorm.Model
 Number string
 UserID uint
}

type User struct {
 gorm.Model
 Name string
 CreditCard CreditCard
}

db.Create(&User{
 Name: "jinzhu",
 CreditCard: CreditCard{Number: "411111111111"}
})
// INSERT INTO `users` ...
// INSERT INTO `credit_cards` ...
```

You can skip saving associations with `Select`, `Omit`, for example:

```go
db.Omit("CreditCard").Create(&user)

// skip all associations
db.Omit(clause.Associations).Create(&user)
```

### <span id="default_values">Default Values</span>

You can define default values for fields with tag `default`, for example:

```go
type User struct {
 ID int64
 Name string `gorm:"default:galeone"`
 Age int64 `gorm:"default:18"`
}
```

Then the default value *will be used* when inserting into the database for [zero-value](https://tour.golang.org/basics/12) fields

{% note warn %}
**NOTE** Any zero value like `0`, `''`, `false` won't be saved into the database for those fields defined default value, you might want to use pointer type or Scanner/Valuer to avoid this, for example:
{% endnote %}

```go
type User struct {
 gorm.Model
 Name string
 Age *int `gorm:"default:18"`
 Active sql.NullBool `gorm:"default:true"`
}
```

{% note warn %}
**NOTE** You have to setup the `default` tag for fields having default or virtual/generated value in database, if you want to skip a default value definition when migrating, you could use `default:(-)`, for example:
{% endnote %}

```go
type User struct {
 ID string `gorm:"default:uuid_generate_v3()"` // db func
 FirstName string
 LastName string
 Age uint8
 FullName string `gorm:"->;type:GENERATED ALWAYS AS (concat(firstname,' ',lastname));default:(-);"`
}
```

{% note warn %}
**NOTE** **SQLite** doesn't support some records are default values when batch insert.
See [SQLite Insert stmt](https://www.sqlite.org/lang_insert.html). For example:

```go
type Pet struct {
	Name string `gorm:"default:cat"`
}

// In SQLite, this is not supported, so GORM will build a wrong SQL to raise error:
// INSERT INTO `pets` (`name`) VALUES ("dog"),(DEFAULT) RETURNING `name`
db.Create(&[]Pet{{Name: "dog"}, {}})
```
A viable alternative is to assign default value to fields in the hook, e.g.

```go
func (p *Pet) BeforeCreate(tx *gorm.DB) (err error) {
	if p.Name == "" {
 p.Name = "cat"
	}
}
```

You can see more info in [issues#6335](https://github.com/go-gorm/gorm/issues/6335)
{% endnote %}

When using virtual/generated value, you might need to disable its creating/updating permission, check out [Field-Level Permission](models.html#field_permission)

### <span id="upsert">Upsert / On Conflict</span>

GORM provides compatible Upsert support for different databases

```go
import "gorm.io/gorm/clause"

// Do nothing on conflict
db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

// Update columns to default value on `id` conflict
db.Clauses(clause.OnConflict{
 Columns: []clause.Column{{Name: "id"}},
 DoUpdates: clause.Assignments(map[string]interface{}{"role": "user"}),
}).Create(&users)
// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET ***; SQL Server
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE ***; MySQL

// Use SQL expression
db.Clauses(clause.OnConflict{
 Columns: []clause.Column{{Name: "id"}},
 DoUpdates: clause.Assignments(map[string]interface{}{"count": gorm.Expr("GREATEST(count, VALUES(count))")}),
}).Create(&users)
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `count`=GREATEST(count, VALUES(count));

// Update columns to new value on `id` conflict
db.Clauses(clause.OnConflict{
 Columns: []clause.Column{{Name: "id"}},
 DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
}).Create(&users)
// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET "name"="excluded"."name"; SQL Server
// INSERT INTO "users" *** ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name", "age"="excluded"."age"; PostgreSQL
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `name`=VALUES(name),`age`=VALUES(age); MySQL

// Update all columns to new value on conflict except primary keys and those columns having default values from sql func
db.Clauses(clause.OnConflict{
 UpdateAll: true,
}).Create(&users)
// INSERT INTO "users" *** ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name", "age"="excluded"."age", ...;
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `name`=VALUES(name),`age`=VALUES(age), ...; MySQL
```

Also checkout `FirstOrInit`, `FirstOrCreate` on [Advanced Query](advanced_query.html)

Checkout [Raw SQL and SQL Builder](sql_builder.html) for more details

---
title: Customize Data Types
layout: page
---

GORM provides few interfaces that allow users to define well-supported customized data types for GORM, takes [json](https://github.com/go-gorm/datatypes/blob/master/json.go) as an example

## Implements Customized Data Type

### Scanner / Valuer

The customized data type has to implement the [Scanner](https://pkg.go.dev/database/sql#Scanner) and [Valuer](https://pkg.go.dev/database/sql/driver#Valuer) interfaces, so GORM knowns to how to receive/save it into the database

For example:

```go
type JSON json.RawMessage

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
 bytes, ok := value.([]byte)
 if !ok {
 return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
 }

 result := json.RawMessage{}
 err := json.Unmarshal(bytes, &result)
 *j = JSON(result)
 return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
 if len(j) == 0 {
 return nil, nil
 }
 return json.RawMessage(j).MarshalJSON()
}
```

There are many third party packages implement the `Scanner`/`Valuer` interface, which can be used with GORM together, for example:

```go
import (
 "github.com/google/uuid"
 "github.com/lib/pq"
)

type Post struct {
 ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
 Title string
 Tags pq.StringArray `gorm:"type:text[]"`
}
```

### GormDataTypeInterface

GORM will read column's database type from [tag](models.html#tags) `type`, if not found, will check if the struct implemented interface `GormDBDataTypeInterface` or `GormDataTypeInterface` and will use its result as data type

```go
type GormDataTypeInterface interface {
 GormDataType() string
}

type GormDBDataTypeInterface interface {
 GormDBDataType(*gorm.DB, *schema.Field) string
}
```

The result of `GormDataType` will be used as the general data type and can be obtained from `schema.Field`'s field `DataType`, which might be helpful when [writing plugins](write_plugins.html) or [hooks](hooks.html) for example:

```go
func (JSON) GormDataType() string {
 return "json"
}

type User struct {
 Attrs JSON
}

func (user User) BeforeCreate(tx *gorm.DB) {
 field := tx.Statement.Schema.LookUpField("Attrs")
 if field.DataType == "json" {
 // do something
 }
}
```

`GormDBDataType` usually returns the right data type for current driver when migrating, for example:

```go
func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
 // use field.Tag, field.TagSettings gets field's tags
 // checkout https://github.com/go-gorm/gorm/blob/master/schema/field.go for all options

 // returns different database type based on driver name
 switch db.Dialector.Name() {
 case "mysql", "sqlite":
 return "JSON"
 case "postgres":
 return "JSONB"
 }
 return ""
}
```

If the struct hasn't implemented the `GormDBDataTypeInterface` or `GormDataTypeInterface` interface, GORM will guess its data type from the struct's first field, for example, will use `string` for `NullString`

```go
type NullString struct {
 String string // use the first field's data type
 Valid bool
}

type User struct {
 Name NullString // data type will be string
}
```

### <span id="gorm_valuer_interface">GormValuerInterface</span>

GORM provides a `GormValuerInterface` interface, which can allow to create/update from SQL Expr or value based on context, for example:

```go
// GORM Valuer interface
type GormValuerInterface interface {
 GormValue(ctx context.Context, db *gorm.DB) clause.Expr
}
```

#### Create/Update from SQL Expr

```go
type Location struct {
	X, Y int
}

func (loc Location) GormDataType() string {
 return "geometry"
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
 return clause.Expr{
 SQL: "ST_PointFromText(?)",
 Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y)},
 }
}

// Scan implements the sql.Scanner interface
func (loc *Location) Scan(v interface{}) error {
 // Scan a value into struct from database driver
}

type User struct {
 ID int
 Name string
 Location Location
}

db.Create(&User{
 Name: "jinzhu",
 Location: Location{X: 100, Y: 100},
})
// INSERT INTO `users` (`name`,`point`) VALUES ("jinzhu",ST_PointFromText("POINT(100 100)"))

db.Model(&User{ID: 1}).Updates(User{
 Name: "jinzhu",
 Location: Location{X: 100, Y: 100},
})
// UPDATE `user_with_points` SET `name`="jinzhu",`location`=ST_PointFromText("POINT(100 100)") WHERE `id` = 1
```

You can also create/update with SQL Expr from map, checkout [Create From SQL Expr](create.html#create_from_sql_expr) and [Update with SQL Expression](update.html#update_from_sql_expr) for details

#### Value based on Context

If you want to create or update a value depends on current context, you can also implements the `GormValuerInterface` interface, for example:

```go
type EncryptedString struct {
 Value string
}

func (es EncryptedString) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) {
 if encryptionKey, ok := ctx.Value("TenantEncryptionKey").(string); ok {
 return clause.Expr{SQL: "?", Vars: []interface{}{Encrypt(es.Value, encryptionKey)}}
 } else {
 db.AddError(errors.New("invalid encryption key"))
 }

 return
}
```

### Clause Expression

If you want to build some query helpers, you can make a struct that implements the `clause.Expression` interface:

```go
type Expression interface {
	Build(builder Builder)
}
```

Checkout [JSON](https://github.com/go-gorm/datatypes/blob/master/json.go) and [SQL Builder](sql_builder.html#clauses) for details, the following is an example of usage:

```go
// Generates SQL with clause Expression
db.Find(&user, datatypes.JSONQuery("attributes").HasKey("role"))
db.Find(&user, datatypes.JSONQuery("attributes").HasKey("orgs", "orga"))

// MySQL
// SELECT * FROM `users` WHERE JSON_EXTRACT(`attributes`, '$.role') IS NOT NULL
// SELECT * FROM `users` WHERE JSON_EXTRACT(`attributes`, '$.orgs.orga') IS NOT NULL

// PostgreSQL
// SELECT * FROM "user" WHERE "attributes"::jsonb ? 'role'
// SELECT * FROM "user" WHERE "attributes"::jsonb -> 'orgs' ? 'orga'

db.Find(&user, datatypes.JSONQuery("attributes").Equals("jinzhu", "name"))
// MySQL
// SELECT * FROM `user` WHERE JSON_EXTRACT(`attributes`, '$.name') = "jinzhu"

// PostgreSQL
// SELECT * FROM "user" WHERE json_extract_path_text("attributes"::json,'name') = 'jinzhu'
```

## Customized Data Types Collections

We created a Github repo for customized data types collections [https://github.com/go-gorm/datatypes](https://github.com/go-gorm/datatypes), pull request welcome ;)

---
title: DBResolver
layout: page
---

DBResolver adds multiple databases support to GORM, the following features are supported:

* Multiple sources, replicas
* Read/Write Splitting
* Automatic connection switching based on the working table/struct
* Manual connection switching
* Sources/Replicas load balancing
* Works for RAW SQL
* Transaction

https://github.com/go-gorm/dbresolver

## Usage

```go
import (
 "gorm.io/gorm"
 "gorm.io/plugin/dbresolver"
 "gorm.io/driver/mysql"
)

db, err := gorm.Open(mysql.Open("db1_dsn"), &gorm.Config{})

db.Use(dbresolver.Register(dbresolver.Config{
 // use `db2` as sources, `db3`, `db4` as replicas
 Sources: []gorm.Dialector{mysql.Open("db2_dsn")},
 Replicas: []gorm.Dialector{mysql.Open("db3_dsn"), mysql.Open("db4_dsn")},
 // sources/replicas load balancing policy
 Policy: dbresolver.RandomPolicy{},
 // print sources/replicas mode in logger
 TraceResolverMode: true,
}).Register(dbresolver.Config{
 // use `db1` as sources (DB's default connection), `db5` as replicas for `User`, `Address`
 Replicas: []gorm.Dialector{mysql.Open("db5_dsn")},
}, &User{}, &Address{}).Register(dbresolver.Config{
 // use `db6`, `db7` as sources, `db8` as replicas for `orders`, `Product`
 Sources: []gorm.Dialector{mysql.Open("db6_dsn"), mysql.Open("db7_dsn")},
 Replicas: []gorm.Dialector{mysql.Open("db8_dsn")},
}, "orders", &Product{}, "secondary"))
```

## Automatic connection switching

DBResolver will automatically switch connection based on the working table/struct

For RAW SQL, DBResolver will extract the table name from the SQL to match the resolver, and will use `sources` unless the SQL begins with `SELECT` (excepts `SELECT... FOR UPDATE`), for example:

```go
// `User` Resolver Examples
db.Table("users").Rows() // replicas `db5`
db.Model(&User{}).Find(&AdvancedUser{}) // replicas `db5`
db.Exec("update users set name = ?", "jinzhu") // sources `db1`
db.Raw("select name from users").Row().Scan(&name) // replicas `db5`
db.Create(&user) // sources `db1`
db.Delete(&User{}, "name = ?", "jinzhu") // sources `db1`
db.Table("users").Update("name", "jinzhu") // sources `db1`

// Global Resolver Examples
db.Find(&Pet{}) // replicas `db3`/`db4`
db.Save(&Pet{}) // sources `db2`

// Orders Resolver Examples
db.Find(&Order{}) // replicas `db8`
db.Table("orders").Find(&Report{}) // replicas `db8`
```

## Read/Write Splitting

Read/Write splitting with DBResolver based on the current used [GORM callbacks](https://gorm.io/docs/write_plugins.html).

For `Query`, `Row` callback, will use `replicas` unless `Write` mode specified
For `Raw` callback, statements are considered read-only and will use `replicas` if the SQL starts with `SELECT`

## Manual connection switching

```go
// Use Write Mode: read user from sources `db1`
db.Clauses(dbresolver.Write).First(&user)

// Specify Resolver: read user from `secondary`'s replicas: db8
db.Clauses(dbresolver.Use("secondary")).First(&user)

// Specify Resolver and Write Mode: read user from `secondary`'s sources: db6 or db7
db.Clauses(dbresolver.Use("secondary"), dbresolver.Write).First(&user)
```

## Transaction

When using transaction, DBResolver will keep using the transaction and won't switch to sources/replicas based on configuration

But you can specifies which DB to use before starting a transaction, for example:

```go
// Start transaction based on default replicas db
tx := db.Clauses(dbresolver.Read).Begin()

// Start transaction based on default sources db
tx := db.Clauses(dbresolver.Write).Begin()

// Start transaction based on `secondary`'s sources
tx := db.Clauses(dbresolver.Use("secondary"), dbresolver.Write).Begin()
```

## Load Balancing

GORM supports load balancing sources/replicas based on policy, the policy should be a struct implements following interface:

```go
type Policy interface {
	Resolve([]gorm.ConnPool) gorm.ConnPool
}
```

Currently only the `RandomPolicy` implemented and it is the default option if no other policy specified.

## Connection Pool

```go
db.Use(
 dbresolver.Register(dbresolver.Config{ /* xxx */ }).
 SetConnMaxIdleTime(time.Hour).
 SetConnMaxLifetime(24 * time.Hour).
 SetMaxIdleConns(100).
 SetMaxOpenConns(200)
)
```

---
title: Delete
layout: page
---

## Delete a Record

When deleting a record, the deleted value needs to have primary key or it will trigger a [Batch Delete](#batch_delete), for example:

### Generics API

```go
ctx := context.Background()

// Delete by ID
err := gorm.G[Email](db).Where("id = ?", 10).Delete(ctx)
// DELETE from emails where id = 10;

// Delete with additional conditions
err := gorm.G[Email](db).Where("id = ? AND name = ?", 10, "jinzhu").Delete(ctx)
// DELETE from emails where id = 10 AND name = "jinzhu";
```

### Traditional API

```go
// Email's ID is `10`
db.Delete(&email)
// DELETE from emails where id = 10;

// Delete with additional conditions
db.Where("name = ?", "jinzhu").Delete(&email)
// DELETE from emails where id = 10 AND name = "jinzhu";
```

## Delete with primary key

GORM allows to delete objects using primary key(s) with inline condition, it works with numbers, check out [Query Inline Conditions](query.html#inline_conditions) for details

```go
db.Delete(&User{}, 10)
// DELETE FROM users WHERE id = 10;

db.Delete(&User{}, "10")
// DELETE FROM users WHERE id = 10;

db.Delete(&users, []int{1,2,3})
// DELETE FROM users WHERE id IN (1,2,3);
```

## Delete Hooks

GORM allows hooks `BeforeDelete`, `AfterDelete`, those methods will be called when deleting a record, refer [Hooks](hooks.html) for details

```go
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Role == "admin" {
 return errors.New("admin user not allowed to delete")
	}
	return
}
```

## <span id="batch_delete">Batch Delete</span>

The specified value has no primary value, GORM will perform a batch delete, it will delete all matched records

### Generics API

```go
ctx := context.Background()

// Batch delete with conditions
err := gorm.G[Email](db).Where("email LIKE ?", "%jinzhu%").Delete(ctx)
// DELETE from emails where email LIKE "%jinzhu%";
```

### Traditional API

```go
db.Where("email LIKE ?", "%jinzhu%").Delete(&Email{})
// DELETE from emails where email LIKE "%jinzhu%";

db.Delete(&Email{}, "email LIKE ?", "%jinzhu%")
// DELETE from emails where email LIKE "%jinzhu%";
```

To efficiently delete large number of records, pass a slice with primary keys to the `Delete` method.

```go
var users = []User{{ID: 1}, {ID: 2}, {ID: 3}}
db.Delete(&users)
// DELETE FROM users WHERE id IN (1,2,3);

db.Delete(&users, "name LIKE ?", "%jinzhu%")
// DELETE FROM users WHERE name LIKE "%jinzhu%" AND id IN (1,2,3);
```

### Block Global Delete

If you perform a batch delete without any conditions, GORM WON'T run it, and will return `ErrMissingWhereClause` error

You have to use some conditions or use raw SQL or enable `AllowGlobalUpdate` mode, for example:

#### Generics API

```go
ctx := context.Background()

// These will return error
err := gorm.G[User](db).Delete(ctx) // gorm.ErrMissingWhereClause

// These will work
err := gorm.G[User](db).Where("1 = 1").Delete(ctx)
// DELETE FROM `users` WHERE 1=1
```

#### Traditional API

```go
db.Delete(&User{}).Error // gorm.ErrMissingWhereClause

db.Delete(&[]User{{Name: "jinzhu1"}, {Name: "jinzhu2"}}).Error // gorm.ErrMissingWhereClause

db.Where("1 = 1").Delete(&User{})
// DELETE FROM `users` WHERE 1=1

db.Exec("DELETE FROM users")
// DELETE FROM users

db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})
// DELETE FROM users
```

### Returning Data From Deleted Rows

Return deleted data, only works for database support Returning, for example:

```go
// return all columns
var users []User
DB.Clauses(clause.Returning{}).Where("role = ?", "admin").Delete(&users)
// DELETE FROM `users` WHERE role = "admin" RETURNING *
// users => []User{{ID: 1, Name: "jinzhu", Role: "admin", Salary: 100}, {ID: 2, Name: "jinzhu.2", Role: "admin", Salary: 1000}}

// return specified columns
DB.Clauses(clause.Returning{Columns: []clause.Column{{Name: "name"}, {Name: "salary"}}}).Where("role = ?", "admin").Delete(&users)
// DELETE FROM `users` WHERE role = "admin" RETURNING `name`, `salary`
// users => []User{{ID: 0, Name: "jinzhu", Role: "", Salary: 100}, {ID: 0, Name: "jinzhu.2", Role: "", Salary: 1000}}
```

## Soft Delete

If your model includes a `gorm.DeletedAt` field (which is included in `gorm.Model`), it will get soft delete ability automatically!

When calling `Delete`, the record WON'T be removed from the database, but GORM will set the `DeletedAt`'s value to the current time, and the data is not findable with normal Query methods anymore.

```go
// user's ID is `111`
db.Delete(&user)
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;

// Batch Delete
db.Where("age = ?", 20).Delete(&User{})
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

// Soft deleted records will be ignored when querying
db.Where("age = 20").Find(&user)
// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
```

If you don't want to include `gorm.Model`, you can enable the soft delete feature like:

```go
type User struct {
 ID int
 Deleted gorm.DeletedAt
 Name string
}
```

### Find soft deleted records

You can find soft deleted records with `Unscoped`

```go
db.Unscoped().Where("age = 20").Find(&users)
// SELECT * FROM users WHERE age = 20;
```

### Delete permanently

You can delete matched records permanently with `Unscoped`

```go
db.Unscoped().Delete(&order)
// DELETE FROM orders WHERE id=10;
```

### Delete Flag

By default, `gorm.Model` uses `*time.Time` as the value for the `DeletedAt` field, and it provides other data formats support with plugin `gorm.io/plugin/soft_delete`

{% note warn %}
**INFO** when creating unique composite index for the DeletedAt field, you must use other data format like unix second/flag with plugin `gorm.io/plugin/soft_delete`'s help, e.g:

```go
import "gorm.io/plugin/soft_delete"

type User struct {
 ID uint
 Name string `gorm:"uniqueIndex:udx_name"`
 DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:udx_name"`
}
```
{% endnote %}

#### Unix Second

Use unix second as delete flag

```go
import "gorm.io/plugin/soft_delete"

type User struct {
 ID uint
 Name string
 DeletedAt soft_delete.DeletedAt
}

// Query
SELECT * FROM users WHERE deleted_at = 0;

// Delete
UPDATE users SET deleted_at = /* current unix second */ WHERE ID = 1;
```

You can also specify to use `milli` or `nano` seconds as the value, for example:

```go
type User struct {
 ID uint
 Name string
 DeletedAt soft_delete.DeletedAt `gorm:"softDelete:milli"`
 // DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano"`
}

// Query
SELECT * FROM users WHERE deleted_at = 0;

// Delete
UPDATE users SET deleted_at = /* current unix milli second or nano second */ WHERE ID = 1;
```

#### Use `1` / `0` AS Delete Flag

```go
import "gorm.io/plugin/soft_delete"

type User struct {
 ID uint
 Name string
 IsDel soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

// Query
SELECT * FROM users WHERE is_del = 0;

// Delete
UPDATE users SET is_del = 1 WHERE ID = 1;
```

#### Mixed Mode

Mixed mode can use `0`, `1` or unix seconds to mark data as deleted or not, and save the deleted time at the same time.

```go
type User struct {
 ID uint
 Name string
 DeletedAt time.Time
 IsDel soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
 // IsDel soft_delete.DeletedAt `gorm:"softDelete:,DeletedAtField:DeletedAt"` // use `unix second`
 // IsDel soft_delete.DeletedAt `gorm:"softDelete:nano,DeletedAtField:DeletedAt"` // use `unix nano second`
}

// Query
SELECT * FROM users WHERE is_del = 0;

// Delete
UPDATE users SET is_del = 1, deleted_at = /* current unix second */ WHERE ID = 1;
```

---
title: Error Handling
layout: page
---

Effective error handling is a cornerstone of robust application development in Go, particularly when interacting with databases using GORM. GORM's approach to error handling requires a nuanced understanding based on the API style you're using.

## Basic Error Handling

### Generics API

With the Generics API, errors are returned directly from the operation methods, following Go's standard error handling pattern:

```go
ctx := context.Background()

// Error handling with direct return values
user, err := gorm.G[User](db).Where("name = ?", "jinzhu").First(ctx)
if err != nil {
 // Handle error...
}

// For operations that don't return a result
err := gorm.G[User](db).Where("id = ?", 1).Delete(ctx)
if err != nil {
 // Handle error...
}
```

### Traditional API

With the Traditional API, GORM integrates error handling into its chainable method syntax. The `*gorm.DB` instance contains an `Error` field, which is set when an error occurs. The common practice is to check this field after executing database operations, especially after [Finisher Methods](method_chaining.html#finisher_method).

After a chain of methods, it's crucial to check the `Error` field:

```go
if err := db.Where("name = ?", "jinzhu").First(&user).Error; err != nil {
 // Handle error...
}
```

Or alternatively:

```go
if result := db.Where("name = ?", "jinzhu").First(&user); result.Error != nil {
 // Handle error...
}
```

## `ErrRecordNotFound`

GORM returns `ErrRecordNotFound` when no record is found using methods like `First`, `Last`, `Take`.

### Generics API

```go
ctx := context.Background()

user, err := gorm.G[User](db).First(ctx)
if errors.Is(err, gorm.ErrRecordNotFound) {
 // Handle record not found error...
}
```

### Traditional API

```go
err := db.First(&user, 100).Error
if errors.Is(err, gorm.ErrRecordNotFound) {
 // Handle record not found error...
}
```

## Handling Error Codes

Many databases return errors with specific codes, which can be indicative of various issues like constraint violations, connection problems, or syntax errors. Handling these error codes in GORM requires parsing the error returned by the database and extracting the relevant code.

```go
import (
 "github.com/go-sql-driver/mysql"
 "gorm.io/gorm"
)

// ...

result := db.Create(&newRecord)
if result.Error != nil {
 if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok {
 switch mysqlErr.Number {
 case 1062: // MySQL code for duplicate entry
 // Handle duplicate entry
 // Add cases for other specific error codes
 default:
 // Handle other errors
 }
 } else {
 // Handle non-MySQL errors or unknown errors
 }
}
```

## Dialect Translated Errors

GORM can return specific errors related to the database dialect being used, when `TranslateError` is enabled, GORM converts database-specific errors into its own generalized errors.

```go
db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{TranslateError: true})
```

- **ErrDuplicatedKey**

This error occurs when an insert operation violates a unique constraint:

```go
result := db.Create(&newRecord)
if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
 // Handle duplicated key error...
}
```

- **ErrForeignKeyViolated**

This error is encountered when a foreign key constraint is violated:

```go
result := db.Create(&newRecord)
if errors.Is(result.Error, gorm.ErrForeignKeyViolated) {
 // Handle foreign key violation error...
}
```

By enabling `TranslateError`, GORM provides a more unified way of handling errors across different databases, translating database-specific errors into common GORM error types.

## Errors

For a complete list of errors that GORM can return, refer to the [Errors List](https://github.com/go-gorm/gorm/blob/master/errors.go) in GORM's documentation.

---
title: Generic database interface sql.DB
layout: page
---

GORM provides the method `DB` which returns a generic database interface [\*sql.DB](https://pkg.go.dev/database/sql#DB) from the current `*gorm.DB`

```go
// Get generic database object sql.DB to use its functions
sqlDB, err := db.DB()

// Ping
sqlDB.Ping()

// Close
sqlDB.Close()

// Returns database statistics
sqlDB.Stats()
```

{% note warn %}
**NOTE** If the underlying database connection is not a `*sql.DB`, like in a transaction, it will returns error
{% endnote %}

## Connection Pool

```go
// Get generic database object sql.DB to use its functions
sqlDB, err := db.DB()

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns sets the maximum number of open connections to the database.
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
sqlDB.SetConnMaxLifetime(time.Hour)
```

---
title: GORM Config
layout: page
---

GORM provides Config can be used during initialization

```go
type Config struct {
 SkipDefaultTransaction bool
 NamingStrategy schema.Namer
 Logger logger.Interface
 NowFunc func() time.Time
 DryRun bool
 PrepareStmt bool
 DisableNestedTransaction bool
 AllowGlobalUpdate bool
 DisableAutomaticPing bool
 DisableForeignKeyConstraintWhenMigrating bool
}
```

## SkipDefaultTransaction

GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency, you can disable it during initialization if it is not required

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 SkipDefaultTransaction: true,
})
```

## <span id="naming_strategy">NamingStrategy</span>

GORM allows users to change the naming conventions by overriding the default `NamingStrategy` which need to implements interface `Namer`

```go
type Namer interface {
	TableName(table string) string
	SchemaName(table string) string
	ColumnName(table, column string) string
	JoinTableName(table string) string
	RelationshipFKName(Relationship) string
	CheckerName(table, column string) string
	IndexName(table, column string) string
}
```

The default `NamingStrategy` also provides few options, like:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 NamingStrategy: schema.NamingStrategy{
 TablePrefix: "t_", // table name prefix, table for `User` would be `t_users`
 SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
 NoLowerCase: true, // skip the snake_casing of names
 NameReplacer: strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
 },
})
```

## Logger

Allow to change GORM's default logger by overriding this option, refer [Logger](logger.html) for more details

## <span id="now_func">NowFunc</span>

Change the function to be used when creating a new timestamp

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 NowFunc: func() time.Time {
 return time.Now().Local()
 },
})
```

## DryRun

Generate `SQL` without executing, can be used to prepare or test generated SQL, refer [Session](session.html) for details

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 DryRun: false,
})
```

## PrepareStmt

`PreparedStmt` creates a prepared statement when executing any SQL and caches them to speed up future calls, refer [Session](session.html) for details

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 PrepareStmt: false,
})
```

## DisableNestedTransaction

When using `Transaction` method inside a db transaction, GORM will use `SavePoint(savedPointName)`, `RollbackTo(savedPointName)` to give you the nested transaction support, you could disable it by using the `DisableNestedTransaction` option, refer [Session](session.html) for details

## AllowGlobalUpdate

Enable global update/delete, refer [Session](session.html) for details

## DisableAutomaticPing

GORM automatically ping database after initialized to check database availability, disable it by setting it to `true`

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 DisableAutomaticPing: true,
})
```

## DisableForeignKeyConstraintWhenMigrating

GORM creates database foreign key constraints automatically when `AutoMigrate` or `CreateTable`, disable this by setting it to `true`, refer [Migration](migration.html) for details

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 DisableForeignKeyConstraintWhenMigrating: true,
})
```

---
title: Has Many
layout: page
---

## Has Many

A `has many` association sets up a one-to-many connection with another model, unlike `has one`, the owner could have zero or many instances of models.

For example, if your application includes users and credit card, and each user can have many credit cards.

### Declare
```go
// User has many CreditCards, UserID is the foreign key
type User struct {
 gorm.Model
 CreditCards []CreditCard
}

type CreditCard struct {
 gorm.Model
 Number string
 UserID uint
}
```

### Retrieve
```go
// Retrieve user list with eager loading credit cards
func GetAll(db *gorm.DB) ([]User, error) {
 var users []User
 err := db.Model(&User{}).Preload("CreditCards").Find(&users).Error
 return users, err
}
```

## Override Foreign Key

To define a `has many` relationship, a foreign key must exist. The default foreign key's name is the owner's type name plus the name of its primary key field

For example, to define a model that belongs to `User`, the foreign key should be `UserID`.

To use another field as foreign key, you can customize it with a `foreignKey` tag, e.g:

```go
type User struct {
 gorm.Model
 CreditCards []CreditCard `gorm:"foreignKey:UserRefer"`
}

type CreditCard struct {
 gorm.Model
 Number string
 UserRefer uint
}
```

## Override References

GORM usually uses the owner's primary key as the foreign key's value, for the above example, it is the `User`'s `ID`,

When you assign credit cards to a user, GORM will save the user's `ID` into credit cards' `UserID` field.

You are able to change it with tag `references`, e.g:

```go
type User struct {
 gorm.Model
 MemberNumber string
 CreditCards []CreditCard `gorm:"foreignKey:UserNumber;references:MemberNumber"`
}

type CreditCard struct {
 gorm.Model
 Number string
 UserNumber string
}
```

## CRUD with Has Many

Please checkout [Association Mode](associations.html#Association-Mode) for working with has many relations

## Eager Loading

GORM allows eager loading has many associations with `Preload`, refer [Preloading (Eager loading)](preload.html) for details

## Self-Referential Has Many

```go
type User struct {
 gorm.Model
 Name string
 ManagerID *uint
 Team []User `gorm:"foreignkey:ManagerID"`
}
```

## FOREIGN KEY Constraints

You can setup `OnUpdate`, `OnDelete` constraints with tag `constraint`, it will be created when migrating with GORM, for example:

```go
type User struct {
 gorm.Model
 CreditCards []CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreditCard struct {
 gorm.Model
 Number string
 UserID uint
}
```

You are also allowed to delete selected has many associations with `Select` when deleting, checkout [Delete with Select](associations.html#delete_with_select) for details

---
title: Has One
layout: page
---

## Has One

A `has one` association sets up a one-to-one connection with another model, but with somewhat different semantics (and consequences). This association indicates that each instance of a model contains or possesses one instance of another model.

For example, if your application includes users and credit cards, and each user can only have one credit card.

### Declare
```go
// User has one CreditCard, UserID is the foreign key
type User struct {
 gorm.Model
 CreditCard CreditCard
}

type CreditCard struct {
 gorm.Model
 Number string
 UserID uint
}
```

### Retrieve
```go
// Retrieve user list with eager loading credit card
func GetAll(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Model(&User{}).Preload("CreditCard").Find(&users).Error
	return users, err
}
```

## Override Foreign Key

For a `has one` relationship, a foreign key field must also exist, the owner will save the primary key of the model belongs to it into this field.

The field's name is usually generated with `has one` model's type plus its `primary key`, for the above example it is `UserID`.

When you give a credit card to the user, it will save the User's `ID` into its `UserID` field.

If you want to use another field to save the relationship, you can change it with tag `foreignKey`, e.g:

```go
type User struct {
 gorm.Model
 CreditCard CreditCard `gorm:"foreignKey:UserName"`
 // use UserName as foreign key
}

type CreditCard struct {
 gorm.Model
 Number string
 UserName string
}
```

## Override References

By default, the owned entity will save the `has one` model's primary key into a foreign key, you could change to save another field's value, like using `Name` for the below example.

You are able to change it with tag `references`, e.g:

```go
type User struct {
 gorm.Model
 Name string `gorm:"index"`
 CreditCard CreditCard `gorm:"foreignKey:UserName;references:Name"`
}

type CreditCard struct {
 gorm.Model
 Number string
 UserName string
}
```

## CRUD with Has One

Please checkout [Association Mode](associations.html#Association-Mode) for working with `has one` relations

## Eager Loading

GORM allows eager loading `has one` associations with `Preload` or `Joins`, refer [Preloading (Eager loading)](preload.html) for details

## Self-Referential Has One

```go
type User struct {
 gorm.Model
 Name string
 ManagerID *uint
 Manager *User
}
```

## FOREIGN KEY Constraints

You can setup `OnUpdate`, `OnDelete` constraints with tag `constraint`, it will be created when migrating with GORM, for example:

```go
type User struct {
 gorm.Model
 CreditCard CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreditCard struct {
 gorm.Model
 Number string
 UserID uint
}
```

You are also allowed to delete selected has one associations with `Select` when deleting, checkout [Delete with Select](associations.html#delete_with_select) for details

---
title: Hints
layout: page
---

GORM provides optimizer/index/comment hints support

https://github.com/go-gorm/hints

## Optimizer Hints

```go
import "gorm.io/hints"

db.Clauses(hints.New("hint")).Find(&User{})
// SELECT * /*+ hint */ FROM `users`
```

## Index Hints

```go
import "gorm.io/hints"

db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
// SELECT * FROM `users` USE INDEX (`idx_user_name`)

db.Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find(&User{})
// SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)"

db.Clauses(
	hints.ForceIndex("idx_user_name", "idx_user_id").ForOrderBy(),
	hints.IgnoreIndex("idx_user_name").ForGroupBy(),
).Find(&User{})
// SELECT * FROM `users` FORCE INDEX FOR ORDER BY (`idx_user_name`,`idx_user_id`) IGNORE INDEX FOR GROUP BY (`idx_user_name`)"
```

## Comment Hints

```go
import "gorm.io/hints"

db.Clauses(hints.Comment("select", "master")).Find(&User{})
// SELECT /*master*/ * FROM `users`;

db.Clauses(hints.CommentBefore("insert", "node2")).Create(&user)
// /*node2*/ INSERT INTO `users` ...;

db.Clauses(hints.CommentAfter("select", "node2")).Create(&user)
// /*node2*/ INSERT INTO `users` ...;

db.Clauses(hints.CommentAfter("where", "hint")).Find(&User{}, "id = ?", 1)
// SELECT * FROM `users` WHERE id = ? /* hint */
```

---
title: Hooks
layout: page
---

## Object Life Cycle

Hooks are functions that are called before or after creation/querying/updating/deletion.

If you have defined specified methods for a model, it will be called automatically when creating, updating, querying, deleting, and if any callback returns an error, GORM will stop future operations and rollback current transaction.

The type of hook methods should be `func(*gorm.DB) error`

## Hooks

### Creating an object

Available hooks for creating

```go
// begin transaction
BeforeSave
BeforeCreate
// save before associations
// insert into database
// save after associations
AfterCreate
AfterSave
// commit or rollback transaction
```

Code Example:

```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
 u.UUID = uuid.New()

 if !u.IsValid() {
 err = errors.New("can't save invalid data")
 }
 return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
 if u.ID == 1 {
 tx.Model(u).Update("role", "admin")
 }
 return
}
```

{% note warn %}
**NOTE** Save/Delete operations in GORM are running in transactions by default, so changes made in that transaction are not visible until it is committed, if you return any error in your hooks, the change will be rollbacked
{% endnote %}

```go
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
 if !u.IsValid() {
 return errors.New("rollback invalid user")
 }
 return nil
}
```

### Updating an object

Available hooks for updating

```go
// begin transaction
BeforeSave
BeforeUpdate
// save before associations
// update database
// save after associations
AfterUpdate
AfterSave
// commit or rollback transaction
```

Code Example:

```go
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
 if u.readonly() {
 err = errors.New("read only user")
 }
 return
}

// Updating data in same transaction
func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
 if u.Confirmed {
 tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("verfied", true)
 }
 return
}
```

### Deleting an object

Available hooks for deleting

```go
// begin transaction
BeforeDelete
// delete from database
AfterDelete
// commit or rollback transaction
```

Code Example:

```go
// Updating data in same transaction
func (u *User) AfterDelete(tx *gorm.DB) (err error) {
 if u.Confirmed {
 tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
 }
 return
}
```

### Querying an object

Available hooks for querying

```go
// load data from database
// Preloading (eager loading)
AfterFind
```

Code Example:

```go
func (u *User) AfterFind(tx *gorm.DB) (err error) {
 if u.MemberShip == "" {
 u.MemberShip = "user"
 }
 return
}
```

## Modify current operation

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
 // Modify current operation through tx.Statement, e.g:
 tx.Statement.Select("Name", "Age")
 tx.Statement.AddClause(clause.OnConflict{DoNothing: true})

 // tx is new session mode with the `NewDB` option
 // operations based on it will run inside same transaction but without any current conditions
 var role Role
 err := tx.First(&role, "name = ?", user.Role).Error
 // SELECT * FROM roles WHERE name = "admin"
 // ...
 return err
}
```

---
title: GORM Guides
layout: page
---

The fantastic ORM library for Golang aims to be developer friendly.

## Overview

* Full-Featured ORM
* Associations (Has One, Has Many, Belongs To, Many To Many, Polymorphism, Single-table inheritance)
* Hooks (Before/After Create/Save/Update/Delete/Find)
* Eager loading with `Preload`, `Joins`
* Transactions, Nested Transactions, Save Point, RollbackTo to Saved Point
* Context, Prepared Statement Mode, DryRun Mode
* Batch Insert, FindInBatches, Find/Create with Map, CRUD with SQL Expr and Context Valuer
* SQL Builder, Upsert, Locking, Optimizer/Index/Comment Hints, Named Argument, SubQuery
* Composite Primary Key, Indexes, Constraints
* Auto Migrations
* Logger
* Generics API for type-safe queries and operations
* Extendable, flexible plugin API: Database Resolver (multiple databases, read/write splitting) / Prometheus...
* Every feature comes with tests
* Developer Friendly

## Install

```sh
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
```

## Quick Start

### Generics API (>= v1.30.0)

```go
package main

import (
 "context"
 "gorm.io/driver/sqlite"
 "gorm.io/gorm"
)

type Product struct {
 gorm.Model
 Code string
 Price uint
}

func main() {
 db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
 if err != nil {
 panic("failed to connect database")
 }

 ctx := context.Background()

 // Migrate the schema
 db.AutoMigrate(&Product{})

 // Create
 err = gorm.G[Product](db).Create(ctx, &Product{Code: "D42", Price: 100})

 // Read
 product, err := gorm.G[Product](db).Where("id = ?", 1).First(ctx) // find product with integer primary key
 products, err := gorm.G[Product](db).Where("code = ?", "D42").Find(ctx) // find product with code D42

 // Update - update product's price to 200
 err = gorm.G[Product](db).Where("id = ?", product.ID).Update(ctx, "Price", 200)
 // Update - update multiple fields
 err = gorm.G[Product](db).Where("id = ?", product.ID).Updates(ctx, Product{Code: "D42", Price: 100})

 // Delete - delete product
 err = gorm.G[Product](db).Where("id = ?", product.ID).Delete(ctx)
}
```

### Traditional API

```go
package main

import (
 "gorm.io/driver/sqlite"
 "gorm.io/gorm"
)

type Product struct {
 gorm.Model
 Code string
 Price uint
}

func main() {
 db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
 if err != nil {
 panic("failed to connect database")
 }

 // Migrate the schema
 db.AutoMigrate(&Product{})

 // Create
 db.Create(&Product{Code: "D42", Price: 100})

 // Read
 var product Product
 db.First(&product, 1) // find product with integer primary key
 db.First(&product, "code = ?", "D42") // find product with code D42

 // Update - update product's price to 200
 db.Model(&product).Update("Price", 200)
 // Update - update multiple fields
 db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
 db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

 // Delete - delete product
 db.Delete(&product, 1)
}
```

---
title: Database Indexes
layout: page
---

GORM allows create database index with tag `index`, `uniqueIndex`, those indexes will be created when [AutoMigrate or CreateTable with GORM](migration.html)

## Index Tag

GORM accepts lots of index settings, like `class`, `type`, `where`, `comment`, `expression`, `sort`, `collate`, `option`

Check the following example for how to use it

```go
type User struct {
	Name string `gorm:"index"`
	Name2 string `gorm:"index:idx_name,unique"`
	Name3 string `gorm:"index:,sort:desc,collate:utf8,type:btree,length:10,where:name3 != 'jinzhu'"`
	Name4 string `gorm:"uniqueIndex"`
	Age int64 `gorm:"index:,class:FULLTEXT,comment:hello \\, world,where:age > 10"`
	Age2 int64 `gorm:"index:,expression:ABS(age)"`
}

// MySQL option
type User struct {
	Name string `gorm:"index:,class:FULLTEXT,option:WITH PARSER ngram INVISIBLE"`
}

// PostgreSQL option
type User struct {
	Name string `gorm:"index:,option:CONCURRENTLY"`
}
```

### uniqueIndex

tag `uniqueIndex` works similar like `index`, it equals to `index:,unique`

```go
type User struct {
	Name1 string `gorm:"uniqueIndex"`
	Name2 string `gorm:"uniqueIndex:idx_name,sort:desc"`
}
```
Note that this will not work for unique composite indexes.

## Composite Indexes

Use same index name for two fields will creates composite indexes, for example:

```go
// create composite index `idx_member` with columns `name`, `number`
type User struct {
	Name string `gorm:"index:idx_member"`
	Number string `gorm:"index:idx_member"`
}
```

For a unique composite index:

```go
// create unique composite index `idx_member` with columns `name`, and `number`
type User struct {
	Name string `gorm:"index:idx_member,unique"`
	Number string `gorm:"index:idx_member,unique"`
}
```

### Fields Priority

The column order of a composite index has an impact on its performance so it must be chosen carefully

You can specify the order with the `priority` option, the default priority value is `10`, if priority value is the same, the order will be based on model struct's field index

```go
type User struct {
	Name string `gorm:"index:idx_member"`
	Number string `gorm:"index:idx_member"`
}
// column order: name, number

type User struct {
	Name string `gorm:"index:idx_member,priority:2"`
	Number string `gorm:"index:idx_member,priority:1"`
}
// column order: number, name

type User struct {
	Name string `gorm:"index:idx_member,priority:12"`
	Number string `gorm:"index:idx_member"`
}
// column order: number, name
```

### Shared composite indexes

If you are creating shared composite indexes with an embedding struct, you can't specify the index name, as embedding the struct more than once results in the duplicated index name in db.

In this case, you can use index tag `composite`, it means the id of the composite index. All fields which have the same composite id of the struct are put together to the same index, just like the original rule. But the improvement is it lets the most derived/embedding struct generates the name of index by NamingStrategy. For example:

```go
type Foo struct {
 IndexA int `gorm:"index:,unique,composite:myname"`
 IndexB int `gorm:"index:,unique,composite:myname"`
}
```

If the table Foo is created, the name of composite index will be `idx_foo_myname`.

```go
type Bar0 struct {
 Foo
}

type Bar1 struct {
 Foo
}
```

Respectively, the name of composite index is `idx_bar0_myname` and `idx_bar1_myname`.

`composite` only works if not specify the name of index.

## Multiple indexes

A field accepts multiple `index`, `uniqueIndex` tags that will create multiple indexes on a field

```go
type UserIndex struct {
	OID int64 `gorm:"index:idx_id;index:idx_oid,unique"`
	MemberNumber string `gorm:"index:idx_id"`
}
```

---
title: Logger
layout: page
---

## Logger

Gorm has a [default logger implementation](https://github.com/go-gorm/gorm/blob/master/logger/logger.go), it will print Slow SQL and happening errors by default

The logger accepts few options, you can customize it during initialization, for example:

```go
newLogger := logger.New(
 log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
 logger.Config{
 SlowThreshold: time.Second, // Slow SQL threshold
 LogLevel: logger.Silent, // Log level
 IgnoreRecordNotFoundError: true, // Ignore ErrRecordNotFound error for logger
 ParameterizedQueries: true, // Don't include params in the SQL log
 Colorful: false, // Disable color
 },
)

// Globally mode
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
 Logger: newLogger,
})

// Continuous session mode
tx := db.Session(&Session{Logger: newLogger})
tx.First(&user)
tx.Model(&user).Update("Age", 18)
```

### Log Levels

GORM defined log levels: `Silent`, `Error`, `Warn`, `Info`

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
 Logger: logger.Default.LogMode(logger.Silent),
})
```

### Debug

Debug a single operation, change current operation's log level to logger.Info

```go
db.Debug().Where("name = ?", "jinzhu").First(&User{})
```

## Customize Logger

Refer to GORM's [default logger](https://github.com/go-gorm/gorm/blob/master/logger/logger.go) for how to define your own one

The logger needs to implement the following interface, it accepts `context`, so you can use it for log tracing

```go
type Interface interface {
	LogMode(LogLevel) Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}
```

---
title: Many To Many
layout: page
---

## Many To Many

Many to Many add a join table between two models.

For example, if your application includes users and languages, and a user can speak many languages, and many users can speak a specified language.

```go
// User has and belongs to many languages, `user_languages` is the join table
type User struct {
 gorm.Model
 Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
 gorm.Model
 Name string
}
```

When using GORM `AutoMigrate` to create a table for `User`, GORM will create the join table automatically

## Back-Reference

### Declare
```go
// User has and belongs to many languages, use `user_languages` as join table
type User struct {
 gorm.Model
 Languages []*Language `gorm:"many2many:user_languages;"`
}

type Language struct {
 gorm.Model
 Name string
 Users []*User `gorm:"many2many:user_languages;"`
}
```

### Retrieve
```go
// Retrieve user list with eager loading languages
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Model(&User{}).Preload("Languages").Find(&users).Error
	return users, err
}

// Retrieve language list with eager loading users
func GetAllLanguages(db *gorm.DB) ([]Language, error) {
	var languages []Language
	err := db.Model(&Language{}).Preload("Users").Find(&languages).Error
	return languages, err
}
```

## Override Foreign Key

For a `many2many` relationship, the join table owns the foreign key which references two models, for example:

```go
type User struct {
 gorm.Model
 Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
 gorm.Model
 Name string
}

// Join Table: user_languages
// foreign key: user_id, reference: users.id
// foreign key: language_id, reference: languages.id
```

To override them, you can use tag `foreignKey`, `references`, `joinForeignKey`, `joinReferences`, not necessary to use them together, you can just use one of them to override some foreign keys/references

```go
type User struct {
	gorm.Model
	Profiles []Profile `gorm:"many2many:user_profiles;foreignKey:Refer;joinForeignKey:UserReferID;References:UserRefer;joinReferences:ProfileRefer"`
	Refer uint `gorm:"index:,unique"`
}

type Profile struct {
	gorm.Model
	Name string
	UserRefer uint `gorm:"index:,unique"`
}

// Which creates join table: user_profiles
// foreign key: user_refer_id, reference: users.refer
// foreign key: profile_refer, reference: profiles.user_refer
```

{% note warn %}
**NOTE:**
Some databases only allow create database foreign keys that reference on a field having unique index, so you need to specify the `unique index` tag if you are creating database foreign keys when migrating
{% endnote %}

## Self-Referential Many2Many

Self-referencing many2many relationship

```go
type User struct {
 gorm.Model
	Friends []*User `gorm:"many2many:user_friends"`
}

// Which creates join table: user_friends
// foreign key: user_id, reference: users.id
// foreign key: friend_id, reference: users.id
```

## Eager Loading

GORM allows eager loading has many associations with `Preload`, refer [Preloading (Eager loading)](preload.html) for details

## CRUD with Many2Many

Please checkout [Association Mode](associations.html#Association-Mode) for working with many2many relations

## Customize JoinTable

`JoinTable` can be a full-featured model, like having `Soft Delete`，`Hooks` supports and more fields, you can set it up with `SetupJoinTable`, for example:

{% note warn %}
**NOTE:**
Customized join table's foreign keys required to be composited primary keys or composited unique index
{% endnote %}

```go
type Person struct {
 ID int
 Name string
 Addresses []Address `gorm:"many2many:person_addresses;"`
}

type Address struct {
 ID uint
 Name string
}

type PersonAddress struct {
 PersonID int `gorm:"primaryKey"`
 AddressID int `gorm:"primaryKey"`
 CreatedAt time.Time
 DeletedAt gorm.DeletedAt
}

func (PersonAddress) BeforeCreate(db *gorm.DB) error {
 // ...
}

// Change model Person's field Addresses' join table to PersonAddress
// PersonAddress must defined all required foreign keys or it will raise error
err := db.SetupJoinTable(&Person{}, "Addresses", &PersonAddress{})
```

## FOREIGN KEY Constraints

You can setup `OnUpdate`, `OnDelete` constraints with tag `constraint`, it will be created when migrating with GORM, for example:

```go
type User struct {
 gorm.Model
 Languages []Language `gorm:"many2many:user_speaks;"`
}

type Language struct {
 Code string `gorm:"primarykey"`
 Name string
}

// CREATE TABLE `user_speaks` (`user_id` integer,`language_code` text,PRIMARY KEY (`user_id`,`language_code`),CONSTRAINT `fk_user_speaks_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE SET NULL ON UPDATE CASCADE,CONSTRAINT `fk_user_speaks_language` FOREIGN KEY (`language_code`) REFERENCES `languages`(`code`) ON DELETE SET NULL ON UPDATE CASCADE);
```

You are also allowed to delete selected many2many relations with `Select` when deleting, checkout [Delete with Select](associations.html#delete_with_select) for details

## Composite Foreign Keys

If you are using [Composite Primary Keys](composite_primary_key.html) for your models, GORM will enable composite foreign keys by default

You are allowed to override the default foreign keys, to specify multiple foreign keys, just separate those keys' name by commas, for example:

```go
type Tag struct {
 ID uint `gorm:"primaryKey"`
 Locale string `gorm:"primaryKey"`
 Value string
}

type Blog struct {
 ID uint `gorm:"primaryKey"`
 Locale string `gorm:"primaryKey"`
 Subject string
 Body string
 Tags []Tag `gorm:"many2many:blog_tags;"`
 LocaleTags []Tag `gorm:"many2many:locale_blog_tags;ForeignKey:id,locale;References:id"`
 SharedTags []Tag `gorm:"many2many:shared_blog_tags;ForeignKey:id;References:id"`
}

// Join Table: blog_tags
// foreign key: blog_id, reference: blogs.id
// foreign key: blog_locale, reference: blogs.locale
// foreign key: tag_id, reference: tags.id
// foreign key: tag_locale, reference: tags.locale

// Join Table: locale_blog_tags
// foreign key: blog_id, reference: blogs.id
// foreign key: blog_locale, reference: blogs.locale
// foreign key: tag_id, reference: tags.id

// Join Table: shared_blog_tags
// foreign key: blog_id, reference: blogs.id
// foreign key: tag_id, reference: tags.id
```

Also check out [Composite Primary Keys](composite_primary_key.html)

---
title: Method Chaining
layout: page
---

GORM's method chaining feature allows for a smooth and fluent style of coding. Here are examples using both the Traditional API and the Generics API:

### Traditional API

```go
db.Where("name = ?", "jinzhu").Where("age = ?", 18).First(&user)
```

### Generics API (>= v1.30.0)

```go
ctx := context.Background()
user, err := gorm.G[User](db).Where("name = ?", "jinzhu").Where("age = ?", 18).First(ctx)
```

Both APIs support method chaining, but the Generics API provides enhanced type safety and returns errors directly from operation methods.

## Method Categories

GORM organizes methods into three primary categories: `Chain Methods`, `Finisher Methods`, and `New Session Methods`. These categories apply to both the Traditional API and the Generics API.

### Chain Methods

Chain methods are used to modify or append `Clauses` to the current `Statement`. Some common chain methods include:

- `Where`
- `Select`
- `Omit`
- `Joins`
- `Scopes`
- `Preload`
- `Raw` (Note: `Raw` cannot be used in conjunction with other chainable methods to build SQL)

For a comprehensive list, visit [GORM Chainable API](https://github.com/go-gorm/gorm/blob/master/chainable_api.go). Also, the [SQL Builder](sql_builder.html) documentation offers more details about `Clauses`.

### Finisher Methods

Finisher methods are immediate, executing registered callbacks that generate and run SQL commands. This category includes methods:

- `Create`
- `First`
- `Find`
- `Take`
- `Save`
- `Update`
- `Delete`
- `Scan`
- `Row`
- `Rows`

For the full list, refer to [GORM Finisher API](https://github.com/go-gorm/gorm/blob/master/finisher_api.go).

### New Session Methods

GORM defines methods like `Session`, `WithContext`, and `Debug` as New Session Methods, which are essential for creating shareable and reusable `*gorm.DB` instances. For more details, see [Session](session.html) documentation.

## Reusability and Safety

### Traditional API

A critical aspect of GORM's Traditional API is understanding when a `*gorm.DB` instance is safe to reuse. Following a `Chain Method` or `Finisher Method`, GORM returns an initialized `*gorm.DB` instance. This instance is not safe for reuse as it may carry over conditions from previous operations, potentially leading to contaminated SQL queries. For example:

### Example of Unsafe Reuse

```go
queryDB := DB.Where("name = ?", "jinzhu")

// First query
queryDB.Where("age > ?", 10).First(&user)
// SQL: SELECT * FROM users WHERE name = "jinzhu" AND age > 10

// Second query with unintended compounded condition
queryDB.Where("age > ?", 20).First(&user2)
// SQL: SELECT * FROM users WHERE name = "jinzhu" AND age > 10 AND age > 20
```

### Example of Safe Reuse

To safely reuse a `*gorm.DB` instance, use a New Session Method:

```go
queryDB := DB.Where("name = ?", "jinzhu").Session(&gorm.Session{})

// First query
queryDB.Where("age > ?", 10).First(&user)
// SQL: SELECT * FROM users WHERE name = "jinzhu" AND age > 10

// Second query, safely isolated
queryDB.Where("age > ?", 20).First(&user2)
// SQL: SELECT * FROM users WHERE name = "jinzhu" AND age > 20
```

In this scenario, using `Session(&gorm.Session{})` ensures that each query starts with a fresh context, preventing the pollution of SQL queries with conditions from previous operations. This is crucial for maintaining the integrity and accuracy of your database interactions.

### Generics API

One of the significant advantages of GORM's Generics API is that it inherently addresses the SQL pollution issue. With the Generics API, you don't need to worry about reusing instances unsafely because:

1. The context is passed directly to each operation method
2. Errors are returned directly from operation methods
3. The generic interface design prevents condition pollution

Here's an example of how the Generics API handles method chaining safely:

```go
ctx := context.Background()

// Define a reusable query base
genericDB := gorm.G[User](db).Where("name = ?", "jinzhu")

// First query
user1, err1 := genericDB.Where("age > ?", 10).First(ctx)
// SQL: SELECT * FROM users WHERE name = "jinzhu" AND age > 10 LIMIT 1

// Second query, no condition pollution
user2, err2 := genericDB.Where("age > ?", 20).First(ctx)
// SQL: SELECT * FROM users WHERE name = "jinzhu" AND age > 20 LIMIT 1
```

The Generics API design significantly reduces the risk of SQL pollution, making your database interactions more reliable and predictable.

## Examples for Clarity

Let's clarify with a few examples using both the Traditional API and the Generics API:

### Traditional API Examples

- **Example 1: Safe Instance Reuse**

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 'db' is a newly initialized `*gorm.DB`, which is safe to reuse.

db.Where("name = ?", "jinzhu").Where("age = ?", 18).Find(&users)
// The first `Where("name = ?", "jinzhu")` call is a chain method that initializes a `*gorm.DB` instance, or `*gorm.Statement`.
// The second `Where("age = ?", 18)` call adds a new condition to the existing `*gorm.Statement`.
// `Find(&users)` is a finisher method, executing registered Query Callbacks, generating and running:
// SELECT * FROM users WHERE name = 'jinzhu' AND age = 18;

db.Where("name = ?", "jinzhu2").Where("age = ?", 20).Find(&users)
// Here, `Where("name = ?", "jinzhu2")` starts a new chain, creating a fresh `*gorm.Statement`.
// `Where("age = ?", 20)` adds to this new statement.
// `Find(&users)` again finalizes the query, executing and generating:
// SELECT * FROM users WHERE name = 'jinzhu2' AND age = 20;

db.Find(&users)
// Directly calling `Find(&users)` without any `Where` starts a new chain and executes:
// SELECT * FROM users;
```

In this example, each chain of method calls is independent, ensuring clean, non-polluted SQL queries.

- **(Bad) Example 2: Unsafe Instance Reuse**

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 'db' is a newly initialized *gorm.DB, safe for initial reuse.

tx := db.Where("name = ?", "jinzhu")
// `Where("name = ?", "jinzhu")` initializes a `*gorm.Statement` instance, which should not be reused across different logical operations.

// Good case
tx.Where("age = ?", 18).Find(&users)
// Reuses 'tx' correctly for a single logical operation, executing:
// SELECT * FROM users WHERE name = 'jinzhu' AND age = 18

// Bad case
tx.Where("age = ?", 28).Find(&users)
// Incorrectly reuses 'tx', compounding conditions and leading to a polluted query:
// SELECT * FROM users WHERE name = 'jinzhu' AND age = 18 AND age = 28;
```

In this bad example, reusing the `tx` variable leads to compounded conditions, which is generally not desirable.

- **Example 3: Safe Reuse with New Session Methods**

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 'db' is a newly initialized *gorm.DB, safe to reuse.

tx := db.Where("name = ?", "jinzhu").Session(&gorm.Session{})
tx := db.Where("name = ?", "jinzhu").WithContext(context.Background())
tx := db.Where("name = ?", "jinzhu").Debug()
// `Session`, `WithContext`, `Debug` methods return a `*gorm.DB` instance marked as safe for reuse. They base a newly initialized `*gorm.Statement` on the current conditions.

// Good case
tx.Where("age = ?", 18).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' AND age = 18

// Good case
tx.Where("age = ?", 28).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' AND age = 28;
```

In this example, using New Session Methods `Session`, `WithContext`, `Debug` correctly initializes a `*gorm.DB` instance for each logical operation, preventing condition pollution and ensuring each query is distinct and based on the specific conditions provided.

### Generics API Examples

- **Example 4: Method Chaining with Generics API**

```go
ctx := context.Background()

// Initialize a generic DB instance
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

// Chain methods with type safety
user, err := gorm.G[User](db).Where("name = ?", "jinzhu").Where("age = ?", 18).First(ctx)
// SELECT * FROM users WHERE name = 'jinzhu' AND age = 18 LIMIT 1;

// Reuse the generic DB instance safely
genericDB := gorm.G[User](db).Where("name = ?", "jinzhu")

// Multiple operations with the same base conditions
user1, err1 := genericDB.Where("age = ?", 18).First(ctx)
// SELECT * FROM users WHERE name = 'jinzhu' AND age = 18 LIMIT 1;

users, err2 := genericDB.Where("age > ?", 20).Find(ctx)
// SELECT * FROM users WHERE name = 'jinzhu' AND age > 20;

// Raw SQL with type safety
users, err3 := gorm.G[User](db).Raw("SELECT * FROM users WHERE name = ? AND age > ?", "jinzhu", 18).Find(ctx)
```

In this example, the Generics API provides type safety while maintaining the fluent method chaining style. The context is passed directly to the finisher methods (`First`, `Find`), and errors are returned directly from these methods, following Go's standard error handling pattern.

Overall, these examples illustrate the importance of understanding GORM's behavior with respect to method chaining and instance management to ensure accurate and efficient database querying. The Generics API offers a more type-safe and less error-prone approach to method chaining compared to the Traditional API.

---
title: Migration
layout: page
---

## Auto Migration

Automatically migrate your schema, to keep your schema up to date.

{% note warn %}
**NOTE:** AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes. It will change existing column's type if its size, precision changed, or if it's changing from non-nullable to nullable. It **WON'T** delete unused columns to protect your data.
{% endnote %}

```go
db.AutoMigrate(&User{})

db.AutoMigrate(&User{}, &Product{}, &Order{})

// Add table suffix when creating tables
db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
```

{% note warn %}
**NOTE** AutoMigrate creates database foreign key constraints automatically, you can disable this feature during initialization, for example:
{% endnote %}

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 DisableForeignKeyConstraintWhenMigrating: true,
})
```

## Migrator Interface

GORM provides a migrator interface, which contains unified API interfaces for each database that could be used to build your database-independent migrations, for example:

SQLite doesn't support `ALTER COLUMN`, `DROP COLUMN`, GORM will create a new table as the one you are trying to change, copy all data, drop the old table, rename the new table

MySQL doesn't support rename column, index for some versions, GORM will perform different SQL based on the MySQL version you are using

```go
type Migrator interface {
 // AutoMigrate
 AutoMigrate(dst ...interface{}) error

 // Database
 CurrentDatabase() string
 FullDataTypeOf(*schema.Field) clause.Expr

 // Tables
 CreateTable(dst ...interface{}) error
 DropTable(dst ...interface{}) error
 HasTable(dst interface{}) bool
 RenameTable(oldName, newName interface{}) error
 GetTables() (tableList []string, err error)

 // Columns
 AddColumn(dst interface{}, field string) error
 DropColumn(dst interface{}, field string) error
 AlterColumn(dst interface{}, field string) error
 MigrateColumn(dst interface{}, field *schema.Field, columnType ColumnType) error
 HasColumn(dst interface{}, field string) bool
 RenameColumn(dst interface{}, oldName, field string) error
 ColumnTypes(dst interface{}) ([]ColumnType, error)

 // Views
 CreateView(name string, option ViewOption) error
 DropView(name string) error

 // Constraints
 CreateConstraint(dst interface{}, name string) error
 DropConstraint(dst interface{}, name string) error
 HasConstraint(dst interface{}, name string) bool

 // Indexes
 CreateIndex(dst interface{}, name string) error
 DropIndex(dst interface{}, name string) error
 HasIndex(dst interface{}, name string) bool
 RenameIndex(dst interface{}, oldName, newName string) error
}
```

### CurrentDatabase

Returns current using database name

```go
db.Migrator().CurrentDatabase()
```

### Tables

```go
// Create table for `User`
db.Migrator().CreateTable(&User{})

// Append "ENGINE=InnoDB" to the creating table SQL for `User`
db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&User{})

// Check table for `User` exists or not
db.Migrator().HasTable(&User{})
db.Migrator().HasTable("users")

// Drop table if exists (will ignore or delete foreign key constraints when dropping)
db.Migrator().DropTable(&User{})
db.Migrator().DropTable("users")

// Rename old table to new table
db.Migrator().RenameTable(&User{}, &UserInfo{})
db.Migrator().RenameTable("users", "user_infos")
```

### Columns

```go
type User struct {
 Name string
}

// Add name field
db.Migrator().AddColumn(&User{}, "Name")
// Drop name field
db.Migrator().DropColumn(&User{}, "Name")
// Alter name field
db.Migrator().AlterColumn(&User{}, "Name")
// Check column exists
db.Migrator().HasColumn(&User{}, "Name")

type User struct {
 Name string
 NewName string
}

// Rename column to new name
db.Migrator().RenameColumn(&User{}, "Name", "NewName")
db.Migrator().RenameColumn(&User{}, "name", "new_name")

// ColumnTypes
db.Migrator().ColumnTypes(&User{}) ([]gorm.ColumnType, error)

type ColumnType interface {
	Name() string
	DatabaseTypeName() string // varchar
	ColumnType() (columnType string, ok bool) // varchar(64)
	PrimaryKey() (isPrimaryKey bool, ok bool)
	AutoIncrement() (isAutoIncrement bool, ok bool)
	Length() (length int64, ok bool)
	DecimalSize() (precision int64, scale int64, ok bool)
	Nullable() (nullable bool, ok bool)
	Unique() (unique bool, ok bool)
	ScanType() reflect.Type
	Comment() (value string, ok bool)
	DefaultValue() (value string, ok bool)
}
```

### Views

Create views by `ViewOption`. About `ViewOption`:

- `Query` is a [subquery](https://gorm.io/docs/advanced_query.html#SubQuery), which is required.
- If `Replace` is true, exec `CREATE OR REPLACE` otherwise exec `CREATE`.
- If `CheckOption` is not empty, append to sql, e.g. `WITH LOCAL CHECK OPTION`.

{% note warn %}
**NOTE** SQLite currently does not support `Replace` in `ViewOption`
{% endnote %}

```go
query := db.Model(&User{}).Where("age > ?", 20)

// Create View
db.Migrator().CreateView("users_pets", gorm.ViewOption{Query: query})
// CREATE VIEW `users_view` AS SELECT * FROM `users` WHERE age > 20

// Create or Replace View
db.Migrator().CreateView("users_pets", gorm.ViewOption{Query: query, Replace: true})
// CREATE OR REPLACE VIEW `users_pets` AS SELECT * FROM `users` WHERE age > 20

// Create View With Check Option
db.Migrator().CreateView("users_pets", gorm.ViewOption{Query: query, CheckOption: "WITH CHECK OPTION"})
// CREATE VIEW `users_pets` AS SELECT * FROM `users` WHERE age > 20 WITH CHECK OPTION

// Drop View
db.Migrator().DropView("users_pets")
// DROP VIEW IF EXISTS "users_pets"
```

### Constraints

```go
type UserIndex struct {
 Name string `gorm:"check:name_checker,name <> 'jinzhu'"`
}

// Create constraint
db.Migrator().CreateConstraint(&User{}, "name_checker")

// Drop constraint
db.Migrator().DropConstraint(&User{}, "name_checker")

// Check constraint exists
db.Migrator().HasConstraint(&User{}, "name_checker")
```

Create foreign keys for relations

```go
type User struct {
 gorm.Model
 CreditCards []CreditCard
}

type CreditCard struct {
 gorm.Model
 Number string
 UserID uint
}

// create database foreign key for user & credit_cards
db.Migrator().CreateConstraint(&User{}, "CreditCards")
db.Migrator().CreateConstraint(&User{}, "fk_users_credit_cards")
// ALTER TABLE `credit_cards` ADD CONSTRAINT `fk_users_credit_cards` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)

// check database foreign key for user & credit_cards exists or not
db.Migrator().HasConstraint(&User{}, "CreditCards")
db.Migrator().HasConstraint(&User{}, "fk_users_credit_cards")

// drop database foreign key for user & credit_cards
db.Migrator().DropConstraint(&User{}, "CreditCards")
db.Migrator().DropConstraint(&User{}, "fk_users_credit_cards")
```

### Indexes

```go
type User struct {
 gorm.Model
 Name string `gorm:"size:255;index:idx_name,unique"`
}

// Create index for Name field
db.Migrator().CreateIndex(&User{}, "Name")
db.Migrator().CreateIndex(&User{}, "idx_name")

// Drop index for Name field
db.Migrator().DropIndex(&User{}, "Name")
db.Migrator().DropIndex(&User{}, "idx_name")

// Check Index exists
db.Migrator().HasIndex(&User{}, "Name")
db.Migrator().HasIndex(&User{}, "idx_name")

type User struct {
 gorm.Model
 Name string `gorm:"size:255;index:idx_name,unique"`
 Name2 string `gorm:"size:255;index:idx_name_2,unique"`
}
// Rename index name
db.Migrator().RenameIndex(&User{}, "Name", "Name2")
db.Migrator().RenameIndex(&User{}, "idx_name", "idx_name_2")
```

## Constraints

GORM creates constraints when auto migrating or creating table, see [Constraints](constraints.html) or [Database Indexes](indexes.html) for details

## Atlas Integration

[Atlas](https://atlasgo.io) is an open-source database migration tool that has an official integration with GORM.

While GORM's `AutoMigrate` feature works in most cases, at some point you may need to switch to a [versioned migrations](https://atlasgo.io/concepts/declarative-vs-versioned#versioned-migrations) strategy.

Once this happens, the responsibility for planning migration scripts and making sure they are in line with what GORM expects at runtime is moved to developers.

Atlas can automatically plan database schema migrations for developers using the official [GORM Provider](https://github.com/ariga/atlas-provider-gorm). After configuring the provider you can automatically plan migrations by running:
```bash
atlas migrate diff --env gorm
```

To learn how to use Atlas with GORM, check out the [official documentation](https://atlasgo.io/guides/orms/gorm).

## Other Migration Tools

To use GORM with other Go-based migration tools, GORM provides a generic DB interface that might be helpful for you.

```go
// returns `*sql.DB`
db.DB()
```

Refer to [Generic Interface](generic_interface.html) for more details.

---
title: Declaring Models
layout: page
---

GORM simplifies database interactions by mapping Go structs to database tables. Understanding how to declare models in GORM is fundamental for leveraging its full capabilities.

## Declaring Models

Models are defined using normal structs. These structs can contain fields with basic Go types, pointers or aliases of these types, or even custom types, as long as they implement the [Scanner](https://pkg.go.dev/database/sql/?tab=doc#Scanner) and [Valuer](https://pkg.go.dev/database/sql/driver#Valuer) interfaces from the `database/sql` package

Consider the following example of a `User` model:

```go
type User struct {
 ID uint // Standard field for the primary key
 Name string // A regular string field
 Email *string // A pointer to a string, allowing for null values
 Age uint8 // An unsigned 8-bit integer
 Birthday *time.Time // A pointer to time.Time, can be null
 MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
 ActivatedAt sql.NullTime // Uses sql.NullTime for nullable time fields
 CreatedAt time.Time // Automatically managed by GORM for creation time
 UpdatedAt time.Time // Automatically managed by GORM for update time
 ignored string // fields that aren't exported are ignored
}
```

In this model:

- Basic data types like `uint`, `string`, and `uint8` are used directly.
- Pointers to types like `*string` and `*time.Time` indicate nullable fields.
- `sql.NullString` and `sql.NullTime` from the `database/sql` package are used for nullable fields with more control.
- `CreatedAt` and `UpdatedAt` are special fields that GORM automatically populates with the current time when a record is created or updated.
- Non-exported fields (starting with a small letter) are not mapped

In addition to the fundamental features of model declaration in GORM, it's important to highlight the support for serialization through the serializer tag. This feature enhances the flexibility of how data is stored and retrieved from the database, especially for fields that require custom serialization logic, See [Serializer](serializer.html) for a detailed explanation

### Conventions

1. **Primary Key**: GORM uses a field named `ID` as the default primary key for each model.

2. **Table Names**: By default, GORM converts struct names to `snake_case` and pluralizes them for table names. For instance, a `User` struct becomes `users` in the database, and a `GormUserName` becomes `gorm_user_names`.

3. **Column Names**: GORM automatically converts struct field names to `snake_case` for column names in the database.

4. **Timestamp Fields**: GORM uses fields named `CreatedAt` and `UpdatedAt` to automatically track the creation and update times of records.

Following these conventions can greatly reduce the amount of configuration or code you need to write. However, GORM is also flexible, allowing you to customize these settings if the default conventions don't fit your requirements. You can learn more about customizing these conventions in GORM's documentation on [conventions](conventions.html).

### `gorm.Model`

GORM provides a predefined struct named `gorm.Model`, which includes commonly used fields:

```go
// gorm.Model definition
type Model struct {
 ID uint `gorm:"primaryKey"`
 CreatedAt time.Time
 UpdatedAt time.Time
 DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

- **Embedding in Your Struct**: You can embed `gorm.Model` directly in your structs to include these fields automatically. This is useful for maintaining consistency across different models and leveraging GORM's built-in conventions, refer [Embedded Struct](#embedded_struct)

- **Fields Included**:
 - `ID`: A unique identifier for each record (primary key).
 - `CreatedAt`: Automatically set to the current time when a record is created.
 - `UpdatedAt`: Automatically updated to the current time whenever a record is updated.
 - `DeletedAt`: Used for soft deletes (marking records as deleted without actually removing them from the database).

## Advanced

### <span id="field_permission">Field-Level Permission</span>

Exported fields have all permissions when doing CRUD with GORM, and GORM allows you to change the field-level permission with tag, so you can make a field to be read-only, write-only, create-only, update-only or ignored

{% note warn %}
**NOTE** ignored fields won't be created when using GORM Migrator to create table
{% endnote %}

```go
type User struct {
 Name string `gorm:"<-:create"` // allow read and create
 Name string `gorm:"<-:update"` // allow read and update
 Name string `gorm:"<-"` // allow read and write (create and update)
 Name string `gorm:"<-:false"` // allow read, disable write permission
 Name string `gorm:"->"` // readonly (disable write permission unless it configured)
 Name string `gorm:"->;<-:create"` // allow read and create
 Name string `gorm:"->:false;<-:create"` // createonly (disabled read from db)
 Name string `gorm:"-"` // ignore this field when write and read with struct
 Name string `gorm:"-:all"` // ignore this field when write, read and migrate with struct
 Name string `gorm:"-:migration"` // ignore this field when migrate with struct
}
```

### <name id="time_tracking">Creating/Updating Time/Unix (Milli/Nano) Seconds Tracking</span>

GORM use `CreatedAt`, `UpdatedAt` to track creating/updating time by convention, and GORM will set the [current time](gorm_config.html#now_func) when creating/updating if the fields are defined

To use fields with a different name, you can configure those fields with tag `autoCreateTime`, `autoUpdateTime`

If you prefer to save UNIX (milli/nano) seconds instead of time, you can simply change the field's data type from `time.Time` to `int`

```go
type User struct {
 CreatedAt time.Time // Set to current time if it is zero on creating
 UpdatedAt int // Set to current unix seconds on updating or if it is zero on creating
 Updated int64 `gorm:"autoUpdateTime:nano"` // Use unix nano seconds as updating time
 Updated int64 `gorm:"autoUpdateTime:milli"`// Use unix milli seconds as updating time
 Created int64 `gorm:"autoCreateTime"` // Use unix seconds as creating time
}
```

### <span id="embedded_struct">Embedded Struct</span>

For anonymous fields, GORM will include its fields into its parent struct, for example:

```go
type Author struct {
 Name string
 Email string
}

type Blog struct {
 Author
 ID int
 Upvotes int32
}
// equals
type Blog struct {
 ID int64
 Name string
 Email string
 Upvotes int32
}
```

For a normal struct field, you can embed it with the tag `embedded`, for example:

```go
type Author struct {
	Name string
	Email string
}

type Blog struct {
 ID int
 Author Author `gorm:"embedded"`
 Upvotes int32
}
// equals
type Blog struct {
 ID int64
	Name string
	Email string
 Upvotes int32
}
```

And you can use tag `embeddedPrefix` to add prefix to embedded fields' db name, for example:

```go
type Blog struct {
 ID int
 Author Author `gorm:"embedded;embeddedPrefix:author_"`
 Upvotes int32
}
// equals
type Blog struct {
 ID int64
	AuthorName string
	AuthorEmail string
 Upvotes int32
}
```

### <span id="tags">Fields Tags</span>

Tags are optional to use when declaring models, GORM supports the following tags:
Tags are case insensitive, however `camelCase` is preferred. If multiple tags are
used they should be separated by a semicolon (`;`). Characters that have special
meaning to the parser can be escaped with a backslash (`\`) allowing them to be
used as parameter values.

| Tag Name | Description |
| --- | --- |
| column | column db name |
| type | column data type, prefer to use compatible general type, e.g: bool, int, uint, float, string, time, bytes, which works for all databases, and can be used with other tags together, like `not null`, `size`, `autoIncrement`... specified database data type like `varbinary(8)` also supported, when using specified database data type, it needs to be a full database data type, for example: `MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT` |
| serializer | specifies serializer for how to serialize and deserialize data into db, e.g: `serializer:json/gob/unixtime` |
| size | specifies column data size/length, e.g: `size:256` |
| primaryKey | specifies column as primary key |
| unique | specifies column as unique |
| default | specifies column default value |
| precision | specifies column precision |
| scale | specifies column scale |
| not null | specifies column as NOT NULL |
| autoIncrement | specifies column auto incrementable |
| autoIncrementIncrement | auto increment step, controls the interval between successive column values |
| embedded | embed the field |
| embeddedPrefix | column name prefix for embedded fields |
| autoCreateTime | track current time when creating, for `int` fields, it will track unix seconds, use value `nano`/`milli` to track unix nano/milli seconds, e.g: `autoCreateTime:nano` |
| autoUpdateTime | track current time when creating/updating, for `int` fields, it will track unix seconds, use value `nano`/`milli` to track unix nano/milli seconds, e.g: `autoUpdateTime:milli` |
| index | create index with options, use same name for multiple fields creates composite indexes, refer [Indexes](indexes.html) for details |
| uniqueIndex | same as `index`, but create uniqued index |
| check | creates check constraint, eg: `check:age > 13`, refer [Constraints](constraints.html) |
| <- | set field's write permission, `<-:create` create-only field, `<-:update` update-only field, `<-:false` no write permission, `<-` create and update permission |
| -> | set field's read permission, `->:false` no read permission |
| - | ignore this field, `-` no read/write permission, `-:migration` no migrate permission, `-:all` no read/write/migrate permission |
| comment | add comment for field when migration |

### Associations Tags

GORM allows configure foreign keys, constraints, many2many table through tags for Associations, check out the [Associations section](associations.html#tags) for details

---
title: Performance
layout: page
---

GORM optimizes many things to improve the performance, the default performance should be good for most applications, but there are still some tips for how to improve it for your application.

## [Disable Default Transaction](transactions.html)

GORM performs write (create/update/delete) operations inside a transaction to ensure data consistency, which is bad for performance, you can disable it during initialization

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 SkipDefaultTransaction: true,
})
```

## [Caches Prepared Statement](session.html)

Creates a prepared statement when executing any SQL and caches them to speed up future calls

```go
// Globally mode
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 PrepareStmt: true,
})

// Session mode
tx := db.Session(&Session{PrepareStmt: true})
tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)
```

{% note warn %}
**NOTE** Also refer how to enable interpolateparams for MySQL to reduce roundtrip https://github.com/go-sql-driver/mysql#interpolateparams
{% endnote %}

### [SQL Builder with PreparedStmt](sql_builder.html)

Prepared Statement works with RAW SQL also, for example:

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 PrepareStmt: true,
})

db.Raw("select sum(age) from users where role = ?", "admin").Scan(&age)
```

You can also use GORM API to prepare SQL with [DryRun Mode](session.html), and execute it with prepared statement later, checkout [Session Mode](session.html) for details

## Select Fields

By default GORM select all fields when querying, you can use `Select` to specify fields you want

```go
db.Select("Name", "Age").Find(&Users{})
```

Or define a smaller API struct to use the [smart select fields feature](advanced_query.html)

```go
type User struct {
 ID uint
 Name string
 Age int
 Gender string
 // hundreds of fields
}

type APIUser struct {
 ID uint
 Name string
}

// Select `id`, `name` automatically when query
db.Model(&User{}).Limit(10).Find(&APIUser{})
// SELECT `id`, `name` FROM `users` LIMIT 10
```

## [Iteration / FindInBatches](advanced_query.html)

Query and process records with iteration or in batches

## [Index Hints](hints.html)

[Index](indexes.html) is used to speed up data search and SQL query performance. `Index Hints` gives the optimizer information about how to choose indexes during query processing, which gives the flexibility to choose a more efficient execution plan than the optimizer

```go
import "gorm.io/hints"

db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
// SELECT * FROM `users` USE INDEX (`idx_user_name`)

db.Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find(&User{})
// SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)"

db.Clauses(
	hints.ForceIndex("idx_user_name", "idx_user_id").ForOrderBy(),
	hints.IgnoreIndex("idx_user_name").ForGroupBy(),
).Find(&User{})
// SELECT * FROM `users` FORCE INDEX FOR ORDER BY (`idx_user_name`,`idx_user_id`) IGNORE INDEX FOR GROUP BY (`idx_user_name`)"
```

## Read/Write Splitting

Increase data throughput through read/write splitting, check out [Database Resolver](dbresolver.html)

---
title: Polymorphism
layout: Page
---

## Polymorphism Association

GORM supports polymorphism association for `has one` and `has many`, it will save owned entity's table name into polymorphic type's field, primary key value into the polymorphic field

By default `polymorphic:<value>` will prefix the column type and column id with `<value>`.
The value will be the table name pluralized.

```go
type Dog struct {
 ID int
 Name string
 Toys []Toy `gorm:"polymorphic:Owner;"`
}

type Toy struct {
 ID int
 Name string
 OwnerID int
 OwnerType string
}

db.Create(&Dog{Name: "dog1", Toys: []Toy{{Name: "toy1"}, {Name: "toy2"}}})
// INSERT INTO `dogs` (`name`) VALUES ("dog1")
// INSERT INTO `toys` (`name`,`owner_id`,`owner_type`) VALUES ("toy1",1,"dogs"), ("toy2",1,"dogs")
```

You can specify polymorphism properties separately using the following GORM tags:

- `polymorphicType`: Specifies the column type.
- `polymorphicId`: Specifies the column ID.
- `polymorphicValue`: Specifies the value of the type.

```go
type Dog struct {
 ID int
 Name string
 Toys []Toy `gorm:"polymorphicType:Kind;polymorphicId:OwnerID;polymorphicValue:master"`
}

type Toy struct {
 ID int
 Name string
 OwnerID int
 Kind string
}

db.Create(&Dog{Name: "dog1", Toys: []Toy{{Name: "toy1"}, {Name: "toy2"}}})
// INSERT INTO `dogs` (`name`) VALUES ("dog1")
// INSERT INTO `toys` (`name`,`owner_id`,`kind`) VALUES ("toy1",1,"master"), ("toy2",1,"master")
```

In these examples, we've used a has-many relationship, but the same principles apply to has-one relationships.

---
title: Preloading (Eager Loading)
layout: page
---

## Preload

GORM allows eager loading relations in other SQL with `Preload`, for example:

### Generics API

```go
type User struct {
 gorm.Model
 Username string
 Orders []Order
}

type Order struct {
 gorm.Model
 UserID uint
 Price float64
}

// Preload Orders when find users
user, err := gorm.G[User](db).Preload("Order", nil).Find(ctx)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4);

// Custom Preloading SQL
user, err = gorm.G[User](db).Preload("Orders", func(db gorm.PreloadBuilder) error {
 db.Order("orders.price DESC")
 return nil
}).Find(ctx)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4) order by orders.price DESC;

```

### Traditional API

```go
// Preload Orders when find users
db.Preload("Orders").Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4);

db.Preload("Orders").Preload("Profile").Preload("Role").Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4); // has many
// SELECT * FROM profiles WHERE user_id IN (1,2,3,4); // has one
// SELECT * FROM roles WHERE id IN (4,5,6); // belongs to
```

## Joins Preloading

`Preload` loads the association data in a separate query, `Join Preload` will loads association data using left join, for example:

### Generics API

```go
type User struct {
 gorm.Model
 Username string
 Order Order
}

type Order struct {
 gorm.Model
 UserID uint
 Price float64
}

users, err := gorm.G[User](db).Joins(clause.JoinTarget{Association: "Order"}, nil).Find(ctx)
// SELECT `users`.`id`,`users`.`username`,`Order`.`id` AS `Order__id`,`Order`.`user_id` AS `Order__user_id`,`Order`.`price` AS `Order__price` FROM `users` JOIN `orders` `Order` ON `users`.`id` = `Order`.`user_id`

// Custom Preloading SQL
users, err := gorm.G[User](db).Joins(
 clause.JoinTarget{Association: "Order"},
 func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
 db.Where(Order{Price: 1000})
 return nil
}).Find(ctx)
// SELECT `users`.`id`,`users`.`username`,`Order`.`id` AS `Order__id`,`Order`.`user_id` AS `Order__user_id`,`Order`.`price` AS `Order__price` FROM `users` JOIN `orders` `Order` ON `users`.`id` = `Order`.`user_id` AND `Order`.`price` = 1000

```

### Traditional API

```go
db.Joins("Company").Joins("Manager").Joins("Account").First(&user, 1)
db.Joins("Company").Joins("Manager").Joins("Account").First(&user, "users.name = ?", "jinzhu")
db.Joins("Company").Joins("Manager").Joins("Account").Find(&users, "users.id IN ?", []int{1,2,3,4,5})
```

Join with conditions

```go
db.Joins("Company", DB.Where(&Company{Alive: true})).Find(&users)
// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` LEFT JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id` AND `Company`.`alive` = true;
```

Join nested model

```go
db.Joins("Manager").Joins("Manager.Company").Find(&users)
// SELECT "users"."id","users"."created_at","users"."updated_at","users"."deleted_at","users"."name","users"."age","users"."birthday","users"."company_id","users"."manager_id","users"."active","Manager"."id" AS "Manager__id","Manager"."created_at" AS "Manager__created_at","Manager"."updated_at" AS "Manager__updated_at","Manager"."deleted_at" AS "Manager__deleted_at","Manager"."name" AS "Manager__name","Manager"."age" AS "Manager__age","Manager"."birthday" AS "Manager__birthday","Manager"."company_id" AS "Manager__company_id","Manager"."manager_id" AS "Manager__manager_id","Manager"."active" AS "Manager__active","Manager__Company"."id" AS "Manager__Company__id","Manager__Company"."name" AS "Manager__Company__name" FROM "users" LEFT JOIN "users" "Manager" ON "users"."manager_id" = "Manager"."id" AND "Manager"."deleted_at" IS NULL LEFT JOIN "companies" "Manager__Company" ON "Manager"."company_id" = "Manager__Company"."id" WHERE "users"."deleted_at" IS NULL
```

{% note warn %}
**NOTE** `Join Preload` works with one-to-one relation, e.g: `has one`, `belongs to`
{% endnote %}

## Preload All

`clause.Associations` can work with `Preload` similar like `Select` when creating/updating, you can use it to `Preload` all associations, for example:

```go
type User struct {
 gorm.Model
 Name string
 CompanyID uint
 Company Company
 Role Role
 Orders []Order
}

db.Preload(clause.Associations).Find(&users)
```

`clause.Associations` won't preload nested associations, but you can use it with [Nested Preloading](#nested_preloading) together, e.g:

```go
db.Preload("Orders.OrderItems.Product").Preload(clause.Associations).Find(&users)
```

## Preload with conditions

GORM allows Preload associations with conditions, it works similar to [Inline Conditions](query.html#inline_conditions)

```go
// Preload Orders with conditions
db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4) AND state NOT IN ('cancelled');

db.Where("state = ?", "active").Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
// SELECT * FROM users WHERE state = 'active';
// SELECT * FROM orders WHERE user_id IN (1,2) AND state NOT IN ('cancelled');
```

## Custom Preloading SQL

You are able to custom preloading SQL by passing in `func(db *gorm.DB) *gorm.DB`, for example:

```go
db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
 return db.Order("orders.amount DESC")
}).Find(&users)
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4) order by orders.amount DESC;
```

## <span id="nested_preloading">Nested Preloading</span>

GORM supports nested preloading, for example:

```go
db.Preload("Orders.OrderItems.Product").Preload("CreditCard").Find(&users)

// Customize Preload conditions for `Orders`
// And GORM won't preload unmatched order's OrderItems then
db.Preload("Orders", "state = ?", "paid").Preload("Orders.OrderItems").Find(&users)
```

## <span id="embedded_preloading">Embedded Preloading</span>

Embedded Preloading is used for [Embedded Struct](models.html#embedded_struct), especially the
same struct. The syntax for Embedded Preloading is similar to Nested Preloading, they are divided by dot.

For example:

```go
type Address struct {
	CountryID int
	Country Country
}

type Org struct {
	PostalAddress Address `gorm:"embedded;embeddedPrefix:postal_address_"`
	VisitingAddress Address `gorm:"embedded;embeddedPrefix:visiting_address_"`
	Address struct {
 ID int
 Address
	} `gorm:"embedded;embeddedPrefix:nested_address_"`
}

// Only preload Org.Address and Org.Address.Country
db.Preload("Address.Country") // "Address" is has_one, "Country" is belongs_to (nested association)

// Only preload Org.PostalAddress
db.Preload("PostalAddress.Country") // "PostalAddress.Country" is belongs_to (embedded association)

// Only preload Org.NestedAddress
db.Preload("NestedAddress.Address.Country") // "NestedAddress.Address.Country" is belongs_to (embedded association)

// All preloaded include "Address" but exclude "Address.Country", because it won't preload nested associations.
db.Preload(clause.Associations)
```

We can omit embedded part when there is no ambiguity.

```go
type Address struct {
	CountryID int
	Country Country
}

type Org struct {
	Address Address `gorm:"embedded"`
}

db.Preload("Address.Country")
db.Preload("Country") // omit "Address" because there is no ambiguity
```

{% note warn %}
**NOTE** `Embedded Preload` only works with `belongs to` relation.
Values of other relations are the same in database, we can't distinguish them.
{% endnote %}

---
title: Prometheus
layout: page
---

GORM provides Prometheus plugin to collect [DBStats](https://pkg.go.dev/database/sql?tab=doc#DBStats) or user-defined metrics

https://github.com/go-gorm/prometheus

## Usage

```go
import (
 "gorm.io/gorm"
 "gorm.io/driver/sqlite"
 "gorm.io/plugin/prometheus"
)

db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

db.Use(prometheus.New(prometheus.Config{
 DBName: "db1", // use `DBName` as metrics label
 RefreshInterval: 15, // Refresh metrics interval (default 15 seconds)
 PushAddr: "prometheus pusher address", // push metrics if `PushAddr` configured
 StartServer: true, // start http server to expose metrics
 HTTPServerPort: 8080, // configure http server port, default port 8080 (if you have configured multiple instances, only the first `HTTPServerPort` will be used to start server)
 MetricsCollector: []prometheus.MetricsCollector {
 &prometheus.MySQL{
 VariableNames: []string{"Threads_running"},
 },
 }, // user defined metrics
}))
```

## User-Defined Metrics

You can define your metrics and collect them with GORM Prometheus plugin, which needs to implements `MetricsCollector` interface

```go
type MetricsCollector interface {
 Metrics(*Prometheus) []prometheus.Collector
}
```

### MySQL

GORM provides an example for how to collect MySQL Status as metrics, check it out [prometheus.MySQL](https://github.com/go-gorm/prometheus/blob/master/mysql.go)

```go
&prometheus.MySQL{
 Prefix: "gorm_status_",
 // Metrics name prefix, default is `gorm_status_`
 // For example, Threads_running's metric name is `gorm_status_Threads_running`
 Interval: 100,
 // Fetch interval, default use Prometheus's RefreshInterval
 VariableNames: []string{"Threads_running"},
 // Select variables from SHOW STATUS, if not set, uses all status variables
}
```

---
title: Query
layout: page
---

## Retrieving a single object

GORM provides `First`, `Take`, `Last` methods to retrieve a single object from the database, it adds `LIMIT 1` condition when querying the database, and it will return the error `ErrRecordNotFound` if no record is found.

### Generics API

```go
ctx := context.Background()

// Get the first record ordered by primary key
user, err := gorm.G[User](db).First(ctx)
// SELECT * FROM users ORDER BY id LIMIT 1;

// Get one record, no specified order
user, err := gorm.G[User](db).Take(ctx)
// SELECT * FROM users LIMIT 1;

// Get last record, ordered by primary key desc
user, err := gorm.G[User](db).Last(ctx)
// SELECT * FROM users ORDER BY id DESC LIMIT 1;

// check error ErrRecordNotFound
errors.Is(err, gorm.ErrRecordNotFound)
```

### Traditional API

```go
// Get the first record ordered by primary key
db.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1;

// Get one record, no specified order
db.Take(&user)
// SELECT * FROM users LIMIT 1;

// Get last record, ordered by primary key desc
db.Last(&user)
// SELECT * FROM users ORDER BY id DESC LIMIT 1;

result := db.First(&user)
result.RowsAffected // returns count of records found
result.Error // returns error or nil

// check error ErrRecordNotFound
errors.Is(result.Error, gorm.ErrRecordNotFound)
```

{% note warn %}
If you want to avoid the `ErrRecordNotFound` error, you could use `Find` like `db.Limit(1).Find(&user)`, the `Find` method accepts both struct and slice data
{% endnote %}

{% note warn %}
Using `Find` without a limit for single object `db.Find(&user)` will query the full table and return only the first object which is non-deterministic and not performant
{% endnote %}

The `First` and `Last` methods will find the first and last record (respectively) as ordered by primary key. They only work when a pointer to the destination struct is passed to the methods as argument or when the model is specified using `db.Model()`. Additionally, if no primary key is defined for relevant model, then the model will be ordered by the first field. For example:

```go
var user User
var users []User

// works because destination struct is passed in
db.First(&user)
// SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1

// works because model is specified using `db.Model()`
result := map[string]interface{}{}
db.Model(&User{}).First(&result)
// SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1

// doesn't work
result := map[string]interface{}{}
db.Table("users").First(&result)

// works with Take
result := map[string]interface{}{}
db.Table("users").Take(&result)

// First and Last still add primary key ordering for the current table when
// using Table with aliases, joins, or custom projections.
type Result struct {
 ID int
 Name string
}
var joinResult Result
db.Table("user_languages ul").
 Select("l.id, l.name").
 Joins("JOIN languages l ON l.code = ul.language_code").
 Order("l.name DESC").
 First(&joinResult)
// SELECT l.id, l.name FROM user_languages ul JOIN languages l ON l.code = ul.language_code ORDER BY l.name DESC, `ul`.`id` LIMIT 1

// Use Take when you provide the ordering yourself.
db.Table("user_languages ul").
 Select("l.id, l.name").
 Joins("JOIN languages l ON l.code = ul.language_code").
 Order("l.name DESC").
 Take(&joinResult)
// SELECT l.id, l.name FROM user_languages ul JOIN languages l ON l.code = ul.language_code ORDER BY l.name DESC LIMIT 1

// no primary key defined, results will be ordered by first field (i.e., `Code`)
type Language struct {
 Code string
 Name string
}
db.First(&Language{})
// SELECT * FROM `languages` ORDER BY `languages`.`code` LIMIT 1
```

### Retrieving objects with primary key

Objects can be retrieved using primary key by using [Inline Conditions](#inline_conditions) if the primary key is a number. When working with strings, extra care needs to be taken to avoid SQL Injection; check out [Security](security.html) section for details.

#### Generics API

```go
ctx := context.Background()

// Using numeric primary key
user, err := gorm.G[User](db).Where("id = ?", 10).First(ctx)
// SELECT * FROM users WHERE id = 10;

// Using string primary key
user, err := gorm.G[User](db).Where("id = ?", "10").First(ctx)
// SELECT * FROM users WHERE id = 10;

// Using multiple primary keys
users, err := gorm.G[User](db).Where("id IN ?", []int{1,2,3}).Find(ctx)
// SELECT * FROM users WHERE id IN (1,2,3);

// If the primary key is a string (for example, like a uuid)
user, err := gorm.G[User](db).Where("id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a").First(ctx)
// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
```

#### Traditional API

```go
db.First(&user, 10)
// SELECT * FROM users WHERE id = 10;

db.First(&user, "10")
// SELECT * FROM users WHERE id = 10;

db.Find(&users, []int{1,2,3})
// SELECT * FROM users WHERE id IN (1,2,3);
```

If the primary key is a string (for example, like a uuid), the query will be written as follows:

```go
db.First(&user, "id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a")
// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
```

When the destination object has a primary value, the primary key will be used to build the condition, for example:

```go
var user = User{ID: 10}
db.First(&user)
// SELECT * FROM users WHERE id = 10;

var result User
db.Model(User{ID: 10}).First(&result)
// SELECT * FROM users WHERE id = 10;
```

{% note warn %}
**NOTE:** If you use gorm's specific field types like `gorm.DeletedAt`, it will run a different query for retrieving object/s.
{% endnote %}

```go
type User struct {
 ID string `gorm:"primarykey;size:16"`
 Name string `gorm:"size:24"`
 DeletedAt gorm.DeletedAt `gorm:"index"`
}

var user = User{ID: 15}
db.First(&user)
// SELECT * FROM `users` WHERE `users`.`id` = '15' AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1
```

## Retrieving all objects

```go
// Get all records
result := db.Find(&users)
// SELECT * FROM users;

result.RowsAffected // returns found records count, equals `len(users)`
result.Error // returns error
```

## Conditions

### String Conditions

```go
// Get first matched record
db.Where("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

// Get all matched records
db.Where("name <> ?", "jinzhu").Find(&users)
// SELECT * FROM users WHERE name <> 'jinzhu';

// IN
db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
// SELECT * FROM users WHERE name IN ('jinzhu','jinzhu 2');

// LIKE
db.Where("name LIKE ?", "%jin%").Find(&users)
// SELECT * FROM users WHERE name LIKE '%jin%';

// AND
db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;

// Time
db.Where("updated_at > ?", lastWeek).Find(&users)
// SELECT * FROM users WHERE updated_at > '2000-01-01 00:00:00';

// BETWEEN
db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
// SELECT * FROM users WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';
```

{% note warn %}
If the object's primary key has been set, then condition query wouldn't cover the value of primary key but use it as a 'and' condition. For example:

```go
var user = User{ID: 10}
db.Where("id = ?", 20).First(&user)
// SELECT * FROM users WHERE id = 10 and id = 20 ORDER BY id ASC LIMIT 1
```

This query would give `record not found` Error. So set the primary key attribute such as `id` to nil before you want to use the variable such as `user` to get new value from database.
{% endnote %}

### Struct & Map Conditions

```go
// Struct
db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;

// Map
db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;

// Slice of primary keys
db.Where([]int64{20, 21, 22}).Find(&users)
// SELECT * FROM users WHERE id IN (20, 21, 22);
```

{% note warn %}
**NOTE** When querying with struct, GORM will only query with non-zero fields, that means if your field's value is `0`, `''`, `false` or other [zero values](https://tour.golang.org/basics/12), it won't be used to build query conditions, for example:
{% endnote %}

```go
db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu";
```

To include zero values in the query conditions, you can use a map, which will include all key-values as query conditions, for example:

```go
db.Where(map[string]interface{}{"Name": "jinzhu", "Age": 0}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;
```

For more details, see [Specify Struct search fields](#specify_search_fields).

### <span id="specify_search_fields">Specify Struct search fields</span>

When searching with struct, you can specify which particular values from the struct to use in the query conditions by passing in the relevant field name or the dbname to `Where()`, for example:

```go
db.Where(&User{Name: "jinzhu"}, "name", "Age").Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;

db.Where(&User{Name: "jinzhu"}, "Age").Find(&users)
// SELECT * FROM users WHERE age = 0;
```

### <span id="inline_conditions">Inline Condition</span>

Query conditions can be inlined into methods like `First` and `Find` in a similar way to `Where`.

```go
// Get by primary key if it were a non-integer type
db.First(&user, "id = ?", "string_primary_key")
// SELECT * FROM users WHERE id = 'string_primary_key';

// Plain SQL
db.Find(&user, "name = ?", "jinzhu")
// SELECT * FROM users WHERE name = "jinzhu";

db.Find(&users, "name <> ? AND age > ?", "jinzhu", 20)
// SELECT * FROM users WHERE name <> "jinzhu" AND age > 20;

// Struct
db.Find(&users, User{Age: 20})
// SELECT * FROM users WHERE age = 20;

// Map
db.Find(&users, map[string]interface{}{"age": 20})
// SELECT * FROM users WHERE age = 20;
```

### Not Conditions

Build NOT conditions, works similar to `Where`

```go
db.Not("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE NOT name = "jinzhu" ORDER BY id LIMIT 1;

// Not In
db.Not(map[string]interface{}{"name": []string{"jinzhu", "jinzhu 2"}}).Find(&users)
// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");

// Struct
db.Not(User{Name: "jinzhu", Age: 18}).First(&user)
// SELECT * FROM users WHERE name <> "jinzhu" AND age <> 18 ORDER BY id LIMIT 1;

// Not In slice of primary keys
db.Not([]int64{1,2,3}).First(&user)
// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;
```

### Or Conditions

```go
db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';

// Struct
db.Where("name = 'jinzhu'").Or(User{Name: "jinzhu 2", Age: 18}).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);

// Map
db.Where("name = 'jinzhu'").Or(map[string]interface{}{"name": "jinzhu 2", "age": 18}).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);
```

For more complicated SQL queries. please also refer to [Group Conditions in Advanced Query](advanced_query.html#group_conditions).

## Selecting Specific Fields

`Select` allows you to specify the fields that you want to retrieve from database. Otherwise, GORM will select all fields by default.

```go
db.Select("name", "age").Find(&users)
// SELECT name, age FROM users;

db.Select([]string{"name", "age"}).Find(&users)
// SELECT name, age FROM users;

db.Table("users").Select("COALESCE(age,?)", 42).Rows()
// SELECT COALESCE(age,'42') FROM users;
```

Also check out [Smart Select Fields](advanced_query.html#smart_select)

## Order

Specify order when retrieving records from the database

```go
db.Order("age desc, name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;

// Multiple orders
db.Order("age desc").Order("name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;

db.Clauses(clause.OrderBy{
 Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{[]int{1, 2, 3}}, WithoutParentheses: true},
}).Find(&User{})
// SELECT * FROM users ORDER BY FIELD(id,1,2,3)
```

## Limit & Offset

`Limit` specify the max number of records to retrieve
`Offset` specify the number of records to skip before starting to return the records

```go
db.Limit(3).Find(&users)
// SELECT * FROM users LIMIT 3;

// Cancel limit condition with -1
db.Limit(10).Find(&users1).Limit(-1).Find(&users2)
// SELECT * FROM users LIMIT 10; (users1)
// SELECT * FROM users; (users2)

db.Offset(3).Find(&users)
// SELECT * FROM users OFFSET 3;

db.Limit(10).Offset(5).Find(&users)
// SELECT * FROM users OFFSET 5 LIMIT 10;

// Cancel offset condition with -1
db.Offset(10).Find(&users1).Offset(-1).Find(&users2)
// SELECT * FROM users OFFSET 10; (users1)
// SELECT * FROM users; (users2)
```

Refer to [Pagination](scopes.html#pagination) for details on how to make a paginator

## Group By & Having

```go
type result struct {
 Date time.Time
 Total int
}

db.Model(&User{}).Select("name, sum(age) as total").Where("name LIKE ?", "group%").Group("name").First(&result)
// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "group%" GROUP BY `name` LIMIT 1

db.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("name = ?", "group").Find(&result)
// SELECT name, sum(age) as total FROM `users` GROUP BY `name` HAVING name = "group"

rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()
defer rows.Close()
for rows.Next() {
 ...
}

rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
defer rows.Close()
for rows.Next() {
 ...
}

type Result struct {
 Date time.Time
 Total int64
}
db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)
```

## Distinct

Selecting distinct values from the model

```go
db.Distinct("name", "age").Order("name, age desc").Find(&results)
```

`Distinct` works with [`Pluck`](advanced_query.html#pluck) and [`Count`](advanced_query.html#count) too

## Joins

### Generics API

The new GORM generics interface brings enhanced support for association queries (`Joins`), offering more flexible association methods, more expressive query capabilities, and a significantly simplified approach to building complex queries.

- **Joins**: Easily specify different join types (e.g., `InnerJoin`, `LeftJoin`) and customize join conditions based on associations, making complex cross-table queries clearer and more intuitive.

```go
// Load only users who have a company
users, err := gorm.G[User](db).Joins(clause.Has("Company"), nil).Find(ctx)

// Use Left Join with custom filter on joined table
user, err = gorm.G[User](db).Joins(clause.LeftJoin.Association("Company"), func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
 db.Where(map[string]any{"name": company.Name})
 return nil
}).Where(map[string]any{"name": user.Name}).First(ctx)

// Join using a subquery
users, err = gorm.G[User](db).Joins(clause.LeftJoin.AssociationFrom("Company", gorm.G[Company](DB).Select("Name")).As("t"),
 func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
 db.Where("?.name = ?", joinTable, u.Company.Name)
 return nil
 },
).Find(ctx)
```

### Traditional API

Specify Joins conditions

```go
type result struct {
 Name string
 Email string
}

db.Model(&User{}).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&result{})
// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id

rows, err := db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()
for rows.Next() {
 ...
}

db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)

// multiple joins with parameter
db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)
```

#### Joins Preloading

You can use `Joins` eager loading associations with a single SQL, for example:

```go
db.Joins("Company").Find(&users)
// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` LEFT JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id`;

// inner join
db.InnerJoins("Company").Find(&users)
// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` INNER JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id`;
```

Join with conditions

```go
db.Joins("Company", db.Where(&Company{Alive: true})).Find(&users)
// SELECT `users`.`id`,`users`.`name`,`users`.`age`,`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `users` LEFT JOIN `companies` AS `Company` ON `users`.`company_id` = `Company`.`id` AND `Company`.`alive` = true;
```

For more details, please refer to [Preloading (Eager Loading)](preload.html).

#### Joins a Derived Table

You can also use `Joins` to join a derived table.

```go
type User struct {
	Id int
	Age int
}

type Order struct {
	UserId int
	FinishedAt *time.Time
}

query := db.Table("order").Select("MAX(order.finished_at) as latest").Joins("left join user user on order.user_id = user.id").Where("user.age > ?", 18).Group("order.user_id")
db.Model(&Order{}).Joins("join (?) q on order.finished_at = q.latest", query).Scan(&results)
// SELECT `order`.`user_id`,`order`.`finished_at` FROM `order` join (SELECT MAX(order.finished_at) as latest FROM `order` left join user user on order.user_id = user.id WHERE user.age > 18 GROUP BY `order`.`user_id`) q on order.finished_at = q.latest
```

## <span id="scan">Scan</span>

Scanning results into a struct works similarly to the way we use `Find`

```go
type Result struct {
 Name string
 Age int
}

var result Result
db.Table("users").Select("name", "age").Where("name = ?", "Antonio").Scan(&result)

// Raw SQL
db.Raw("SELECT name, age FROM users WHERE name = ?", "Antonio").Scan(&result)
```

---
title: Scopes
layout: page
---

Scopes allow you to re-use commonly used logic, the shared logic needs to be defined as type `func(*gorm.DB) *gorm.DB`

## Query

Scope examples for querying

```go
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
 return db.Where("amount > ?", 1000)
}

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
 return db.Where("pay_mode = ?", "card")
}

func PaidWithCod(db *gorm.DB) *gorm.DB {
 return db.Where("pay_mode = ?", "cod")
}

func OrderStatus(status []string) func (db *gorm.DB) *gorm.DB {
 return func (db *gorm.DB) *gorm.DB {
 return db.Scopes(AmountGreaterThan1000).Where("status IN (?)", status)
 }
}

db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)
// Find all credit card orders and amount greater than 1000

db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&orders)
// Find all COD orders and amount greater than 1000

db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
// Find all paid, shipped orders that amount greater than 1000
```

### <span id="pagination">Pagination</span>

```go
func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
 return func (db *gorm.DB) *gorm.DB {
 q := r.URL.Query()
 page, _ := strconv.Atoi(q.Get("page"))
 if page <= 0 {
 page = 1
 }

 pageSize, _ := strconv.Atoi(q.Get("page_size"))
 switch {
 case pageSize > 100:
 pageSize = 100
 case pageSize <= 0:
 pageSize = 10
 }

 offset := (page - 1) * pageSize
 return db.Offset(offset).Limit(pageSize)
 }
}

db.Scopes(Paginate(r)).Find(&users)
db.Scopes(Paginate(r)).Find(&articles)
```

## Dynamically Table

Use `Scopes` to dynamically set the query Table

```go
func TableOfYear(user *User, year int) func(db *gorm.DB) *gorm.DB {
 return func(db *gorm.DB) *gorm.DB {
 tableName := user.TableName() + strconv.Itoa(year)
 return db.Table(tableName)
 }
}

DB.Scopes(TableOfYear(user, 2019)).Find(&users)
// SELECT * FROM users_2019;

DB.Scopes(TableOfYear(user, 2020)).Find(&users)
// SELECT * FROM users_2020;

// Table form different database
func TableOfOrg(user *User, dbName string) func(db *gorm.DB) *gorm.DB {
 return func(db *gorm.DB) *gorm.DB {
 tableName := dbName + "." + user.TableName()
 return db.Table(tableName)
 }
}

DB.Scopes(TableOfOrg(user, "org1")).Find(&users)
// SELECT * FROM org1.users;

DB.Scopes(TableOfOrg(user, "org2")).Find(&users)
// SELECT * FROM org2.users;
```

## Updates

Scope examples for updating/deleting

```go
func CurOrganization(r *http.Request) func(db *gorm.DB) *gorm.DB {
 return func (db *gorm.DB) *gorm.DB {
 org := r.Query("org")

 if org != "" {
 var organization Organization
 if db.Session(&Session{}).First(&organization, "name = ?", org).Error == nil {
 return db.Where("org_id = ?", organization.ID)
 }
 }

 db.AddError("invalid organization")
 return db
 }
}

db.Model(&article).Scopes(CurOrganization(r)).Update("Name", "name 1")
// UPDATE articles SET name = "name 1" WHERE org_id = 111
db.Scopes(CurOrganization(r)).Delete(&Article{})
// DELETE FROM articles WHERE org_id = 111
```

---
title: Security
layout: page
---

GORM uses the `database/sql`'s argument placeholders to construct the SQL statement, which will automatically escape arguments to avoid SQL injection

{% note warn %}
**NOTE** The SQL from Logger is not fully escaped like the one executed, be careful when copying and executing it in SQL console
{% endnote %}

## Query Condition

User's input should be only used as an argument, for example:

```go
userInput := "jinzhu;drop table users;"

// safe, will be escaped
db.Where("name = ?", userInput).First(&user)

// SQL injection
db.Where(fmt.Sprintf("name = %v", userInput)).First(&user)
```

## Inline Condition

```go
// will be escaped
db.First(&user, "name = ?", userInput)

// SQL injection
db.First(&user, fmt.Sprintf("name = %v", userInput))
```

When retrieving objects with number primary key by user's input, you should check the type of variable.

```go
userInputID := "1=1;drop table users;"
// safe, return error
id, err := strconv.Atoi(userInputID)
if err != nil {
 return err
}
db.First(&user, id)

// SQL injection
db.First(&user, userInputID)
// SELECT * FROM users WHERE 1=1;drop table users;
```

## SQL injection Methods

To support some features, some inputs are not escaped, be careful when using user's input with those methods

```go
db.Select("name; drop table users;").First(&user)
db.Distinct("name; drop table users;").First(&user)

db.Model(&user).Pluck("name; drop table users;", &names)

db.Group("name; drop table users;").First(&user)

db.Group("name").Having("1 = 1;drop table users;").First(&user)

db.Raw("select name from users; drop table users;").First(&user)

db.Exec("select name from users; drop table users;")

db.Order("name; drop table users;").First(&user)

db.Table("users; drop table users;").Find(&users)

db.Delete(&User{}, "id=1; drop table users;")

db.Joins("inner join orders; drop table users;").Find(&users)

db.InnerJoins("inner join orders; drop table users;").Find(&users)

//Despite being parameterized in Exec() function, gorm.Expr is still injectable
db.Exec("UPDATE users SET name = '?' WHERE id = 1", gorm.Expr("alice'; drop table users;-- "))

db.Where("id=1").Not("name = 'alice'; drop table users;").Find(&users)

db.Where("id=1").Or("name = 'alice'; drop table users;").Find(&users)

db.Find(&User{}, "name = 'alice'; drop table users;")

// The following functions can only be injected by blind SQL injection methods
db.First(&users, "2 or 1=1-- ")

db.FirstOrCreate(&users, "2 or 1=1-- ")

db.FirstOrInit(&users, "2 or 1=1-- ")

db.Last(&users, "2 or 1=1-- ")

db.Take(&users, "2 or 1=1-- ")
```

The general rule to avoid SQL injection is don't trust user-submitted data, you can perform whitelist validation to test user input against an existing set of known, approved, and defined input, and when using user's input, only use them as an argument.

---
title: Serializer
layout: page
---

Serializer is an extensible interface that allows to customize how to serialize and deserialize data with database.

GORM provides some default serializers: `json`, `gob`, `unixtime`, here is a quick example of how to use it.

```go
type User struct {
	Name []byte `gorm:"serializer:json"`
	Roles Roles `gorm:"serializer:json"`
	Contracts map[string]interface{} `gorm:"serializer:json"`
	JobInfo Job `gorm:"type:bytes;serializer:gob"`
	CreatedTime int64 `gorm:"serializer:unixtime;type:time"` // store int as datetime into database
}

type Roles []string

type Job struct {
	Title string
	Location string
	IsIntern bool
}

createdAt := time.Date(2020, 1, 1, 0, 8, 0, 0, time.UTC)
data := User{
 Name: []byte("jinzhu"),
 Roles: []string{"admin", "owner"},
 Contracts: map[string]interface{}{"name": "jinzhu", "age": 10},
 CreatedTime: createdAt.Unix(),
 JobInfo: Job{
 Title: "Developer",
 Location: "NY",
 IsIntern: false,
 },
}

DB.Create(&data)
// INSERT INTO `users` (`name`,`roles`,`contracts`,`job_info`,`created_time`) VALUES
// ("\"amluemh1\"","[\"admin\",\"owner\"]","{\"age\":10,\"name\":\"jinzhu\"}",<gob binary>,"2020-01-01 00:08:00")

var result User
DB.First(&result, "id = ?", data.ID)
// result => User{
// Name: []byte("jinzhu"),
// Roles: []string{"admin", "owner"},
// Contracts: map[string]interface{}{"name": "jinzhu", "age": 10},
// CreatedTime: createdAt.Unix(),
// JobInfo: Job{
// Title: "Developer",
// Location: "NY",
// IsIntern: false,
// },
// }

DB.Where(User{Name: []byte("jinzhu")}).Take(&result)
// SELECT * FROM `users` WHERE `users`.`name` = "\"amluemh1\"
```

## Register Serializer

A Serializer needs to implement how to serialize and deserialize data, so it requires to implement the the following interface

```go
import "gorm.io/gorm/schema"

type SerializerInterface interface {
	Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error
	SerializerValuerInterface
}

type SerializerValuerInterface interface {
	Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error)
}
```

For example, the default `JSONSerializer` is implemented like:

```go
// JSONSerializer json serializer
type JSONSerializer struct {
}

// Scan implements serializer interface
func (JSONSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	fieldValue := reflect.New(field.FieldType)

	if dbValue != nil {
 var bytes []byte
 switch v := dbValue.(type) {
 case []byte:
 bytes = v
 case string:
 bytes = []byte(v)
 default:
 return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
 }

 err = json.Unmarshal(bytes, fieldValue.Interface())
	}

	field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
	return
}

// Value implements serializer interface
func (JSONSerializer) Value(ctx context.Context, field *Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}
```

And registered with the following code:

```go
schema.RegisterSerializer("json", JSONSerializer{})
```

After registering a serializer, you can use it with the `serializer` tag, for example:

```go
type User struct {
	Name []byte `gorm:"serializer:json"`
}
```

## Customized Serializer Type

You can use a registered serializer with tags, you are also allowed to create a customized struct that implements the above `SerializerInterface` and use it as a field type directly, for example:

```go
type EncryptedString string

// ctx: contains request-scoped values
// field: the field using the serializer, contains GORM settings, struct tags
// dst: current model value, `user` in the below example
// dbValue: current field's value in database
func (es *EncryptedString) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	switch value := dbValue.(type) {
	case []byte:
 *es = EncryptedString(bytes.TrimPrefix(value, []byte("hello")))
	case string:
 *es = EncryptedString(strings.TrimPrefix(value, "hello"))
	default:
 return fmt.Errorf("unsupported data %#v", dbValue)
	}
	return nil
}

// ctx: contains request-scoped values
// field: the field using the serializer, contains GORM settings, struct tags
// dst: current model value, `user` in the below example
// fieldValue: current field's value of the dst
func (es EncryptedString) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return "hello" + string(es), nil
}

type User struct {
	gorm.Model
	Password EncryptedString
}

data := User{
	Password: EncryptedString("pass"),
}

DB.Create(&data)
// INSERT INTO `serializer_structs` (`password`) VALUES ("hellopass")

var result User
DB.First(&result, "id = ?", data.ID)
// result => User{
// Password: EncryptedString("pass"),
// }

DB.Where(User{Password: EncryptedString("pass")}).Take(&result)
// SELECT * FROM `users` WHERE `users`.`password` = "hellopass"
```

---
title: Session
layout: page
---

GORM provides `Session` method, which is a [`New Session Method`](method_chaining.html), it allows to create a new session mode with configuration:

```go
// Session Configuration
type Session struct {
 DryRun bool
 PrepareStmt bool
 NewDB bool
 Initialized bool
 SkipHooks bool
 SkipDefaultTransaction bool
 DisableNestedTransaction bool
 AllowGlobalUpdate bool
 FullSaveAssociations bool
 QueryFields bool
 Context context.Context
 Logger logger.Interface
 NowFunc func() time.Time
 CreateBatchSize int
}
```

## DryRun

Generates `SQL` without executing. It can be used to prepare or test generated SQL, for example:

```go
// session mode
stmt := db.Session(&Session{DryRun: true}).First(&user, 1).Statement
stmt.SQL.String() //=> SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
stmt.Vars //=> []interface{}{1}

// globally mode with DryRun
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{DryRun: true})

// different databases generate different SQL
stmt := db.Find(&user, 1).Statement
stmt.SQL.String() //=> SELECT * FROM `users` WHERE `id` = $1 // PostgreSQL
stmt.SQL.String() //=> SELECT * FROM `users` WHERE `id` = ? // MySQL
stmt.Vars //=> []interface{}{1}
```

To generate the final SQL, you could use following code:

```go
// NOTE: the SQL is not always safe to execute, GORM only uses it for logs, it might cause SQL injection
db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
// SELECT * FROM `users` WHERE `id` = 1
```

## PrepareStmt

`PreparedStmt` creates prepared statements when executing any SQL and caches them to speed up future calls, for example:

```go
// globally mode, all DB operations will create prepared statements and cache them
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 PrepareStmt: true,
})

// session mode
tx := db.Session(&Session{PrepareStmt: true})
tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)

// returns prepared statements manager
stmtManger, ok := tx.ConnPool.(*PreparedStmtDB)

// close prepared statements for *current session*
stmtManger.Close()

// prepared SQL for *current session*
stmtManger.PreparedSQL // => []string{}

// prepared statements for current database connection pool (all sessions)
stmtManger.Stmts // map[string]*sql.Stmt

for sql, stmt := range stmtManger.Stmts {
 sql // prepared SQL
 stmt // prepared statement
 stmt.Close() // close the prepared statement
}
```

## NewDB

Create a new DB without conditions with option `NewDB`, for example:

```go
tx := db.Where("name = ?", "jinzhu").Session(&gorm.Session{NewDB: true})

tx.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1

tx.First(&user, "id = ?", 10)
// SELECT * FROM users WHERE id = 10 ORDER BY id

// Without option `NewDB`
tx2 := db.Where("name = ?", "jinzhu").Session(&gorm.Session{})
tx2.First(&user)
// SELECT * FROM users WHERE name = "jinzhu" ORDER BY id
```

## Initialized

Create a new initialized DB, which is not Method Chain/Goroutine Safe anymore, refer [Method Chaining](method_chaining.html)

```go
tx := db.Session(&gorm.Session{Initialized: true})
```

## Skip Hooks

If you want to skip `Hooks` methods, you can use the `SkipHooks` session mode, for example:

```go
DB.Session(&gorm.Session{SkipHooks: true}).Create(&user)

DB.Session(&gorm.Session{SkipHooks: true}).Create(&users)

DB.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(users, 100)

DB.Session(&gorm.Session{SkipHooks: true}).Find(&user)

DB.Session(&gorm.Session{SkipHooks: true}).Delete(&user)

DB.Session(&gorm.Session{SkipHooks: true}).Model(User{}).Where("age > ?", 18).Updates(&user)
```

## DisableNestedTransaction

When using `Transaction` method inside a DB transaction, GORM will use `SavePoint(savedPointName)`, `RollbackTo(savedPointName)` to give you the nested transaction support. You can disable it by using the `DisableNestedTransaction` option, for example:

```go
db.Session(&gorm.Session{
 DisableNestedTransaction: true,
}).CreateInBatches(&users, 100)
```

## AllowGlobalUpdate

GORM doesn't allow global update/delete by default, will return `ErrMissingWhereClause` error. You can set this option to true to enable it, for example:

```go
db.Session(&gorm.Session{
 AllowGlobalUpdate: true,
}).Model(&User{}).Update("name", "jinzhu")
// UPDATE users SET `name` = "jinzhu"
```

## FullSaveAssociations

GORM will auto-save associations and its reference using [Upsert](create.html#upsert) when creating/updating a record. If you want to update associations' data, you should use the `FullSaveAssociations` mode, for example:

```go
db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)
// ...
// INSERT INTO "addresses" (address1) VALUES ("Billing Address - Address 1"), ("Shipping Address - Address 1") ON DUPLICATE KEY SET address1=VALUES(address1);
// INSERT INTO "users" (name,billing_address_id,shipping_address_id) VALUES ("jinzhu", 1, 2);
// INSERT INTO "emails" (user_id,email) VALUES (111, "jinzhu@example.com"), (111, "jinzhu-2@example.com") ON DUPLICATE KEY SET email=VALUES(email);
// ...
```

## Context

With the `Context` option, you can set the `Context` for following SQL operations, for example:

```go
timeoutCtx, _ := context.WithTimeout(context.Background(), time.Second)
tx := db.Session(&Session{Context: timeoutCtx})

tx.First(&user) // query with context timeoutCtx
tx.Model(&user).Update("role", "admin") // update with context timeoutCtx
```

GORM also provides shortcut method `WithContext`, here is the definition:

```go
func (db *DB) WithContext(ctx context.Context) *DB {
 return db.Session(&Session{Context: ctx})
}
```

## Logger

Gorm allows customizing built-in logger with the `Logger` option, for example:

```go
newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
 logger.Config{
 SlowThreshold: time.Second,
 LogLevel: logger.Silent,
 Colorful: false,
 })
db.Session(&Session{Logger: newLogger})

db.Session(&Session{Logger: logger.Default.LogMode(logger.Silent)})
```

Checkout [Logger](logger.html) for more details.

## NowFunc

`NowFunc` allows changing the function to get current time of GORM, for example:

```go
db.Session(&Session{
 NowFunc: func() time.Time {
 return time.Now().Local()
 },
})
```

## Debug

`Debug` is a shortcut method to change session's `Logger` to debug mode, here is the definition:

```go
func (db *DB) Debug() (tx *DB) {
 return db.Session(&Session{
 Logger: db.Logger.LogMode(logger.Info),
 })
}
```

## QueryFields

Select by fields

```go
db.Session(&gorm.Session{QueryFields: true}).Find(&user)
// SELECT `users`.`name`, `users`.`age`, ... FROM `users` // with this option
// SELECT * FROM `users` // without this option
```

## CreateBatchSize

Default batch size

```go
users = [5000]User{{Name: "jinzhu", Pets: []Pet{pet1, pet2, pet3}}...}

db.Session(&gorm.Session{CreateBatchSize: 1000}).Create(&users)
// INSERT INTO users xxx (5 batches)
// INSERT INTO pets xxx (15 batches)
```

---
title: Settings
layout: page
---

GORM provides `Set`, `Get`, `InstanceSet`, `InstanceGet` methods allow users pass values to [hooks](hooks.html) or other methods

GORM uses this for some features, like pass creating table options when migrating table.

```go
// Add table suffix when creating tables
db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
```

## Set / Get

Use `Set` / `Get` pass settings to hooks methods, for example:

```go
type User struct {
 gorm.Model
 CreditCard CreditCard
 // ...
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
 myValue, ok := tx.Get("my_value")
 // ok => true
 // myValue => 123
}

type CreditCard struct {
 gorm.Model
 // ...
}

func (card *CreditCard) BeforeCreate(tx *gorm.DB) error {
 myValue, ok := tx.Get("my_value")
 // ok => true
 // myValue => 123
}

myValue := 123
db.Set("my_value", myValue).Create(&User{})
```

## InstanceSet / InstanceGet

Use `InstanceSet` / `InstanceGet` pass settings to current `*Statement`'s hooks methods, for example:

```go
type User struct {
 gorm.Model
 CreditCard CreditCard
 // ...
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
 myValue, ok := tx.InstanceGet("my_value")
 // ok => true
 // myValue => 123
}

type CreditCard struct {
 gorm.Model
 // ...
}

// When creating associations, GORM creates a new `*Statement`, so can't read other instance's settings
func (card *CreditCard) BeforeCreate(tx *gorm.DB) error {
 myValue, ok := tx.InstanceGet("my_value")
 // ok => false
 // myValue => nil
}

myValue := 123
db.InstanceSet("my_value", myValue).Create(&User{})
```

---
title: Sharding
layout: page
---

Sharding plugin using SQL parser and replace for splits large tables into smaller ones, redirects Query into sharding tables. Give you a high performance database access.

https://github.com/go-gorm/sharding

## Features

- Non-intrusive design. Load the plugin, specify the config, and all done.
- Lighting-fast. No network based middlewares, as fast as Go.
- Multiple database (PostgreSQL, MySQL) support.
- Integrated primary key generator (Snowflake, PostgreSQL Sequence, Custom, ...).

## Usage

Config the sharding middleware, register the tables which you want to shard. See [Godoc](https://pkg.go.dev/gorm.io/sharding) for config details.

```go
import (
 "fmt"

 "gorm.io/driver/postgres"
 "gorm.io/gorm"
 "gorm.io/sharding"
)

db, err := gorm.Open(postgres.New(postgres.Config{DSN: "postgres://localhost:5432/sharding-db?sslmode=disable"))

db.Use(sharding.Register(sharding.Config{
 ShardingKey: "user_id",
 NumberOfShards: 64,
 PrimaryKeyGenerator: sharding.PKSnowflake,
}, "orders", Notification{}, AuditLog{}))
// This case for show up give notifications, audit_logs table use same sharding rule.
```

Use the db session as usual. Just note that the query should have the `Sharding Key` when operate sharding tables.

```go
// Gorm create example, this will insert to orders_02
db.Create(&Order{UserID: 2})
// sql: INSERT INTO orders_2 ...

// Show have use Raw SQL to insert, this will insert into orders_03
db.Exec("INSERT INTO orders(user_id) VALUES(?)", int64(3))

// This will throw ErrMissingShardingKey error, because there not have sharding key presented.
db.Create(&Order{Amount: 10, ProductID: 100})
fmt.Println(err)

// Find, this will redirect query to orders_02
var orders []Order
db.Model(&Order{}).Where("user_id", int64(2)).Find(&orders)
fmt.Printf("%#v\n", orders)

// Raw SQL also supported
db.Raw("SELECT * FROM orders WHERE user_id = ?", int64(3)).Scan(&orders)
fmt.Printf("%#v\n", orders)

// This will throw ErrMissingShardingKey error, because WHERE conditions not included sharding key
err = db.Model(&Order{}).Where("product_id", "1").Find(&orders).Error
fmt.Println(err)

// Update and Delete are similar to create and query
db.Exec("UPDATE orders SET product_id = ? WHERE user_id = ?", 2, int64(3))
err = db.Exec("DELETE FROM orders WHERE product_id = 3").Error
fmt.Println(err) // ErrMissingShardingKey
```

The full example is [here](https://github.com/go-gorm/sharding/tree/main/examples).

---
title: SQL Builder
layout: page
---

## Raw SQL

### Generics API

Query Raw SQL with `Scan`

```go
type Result struct {
 ID int
 Name string
 Age int
}

// Scan into a single result
result, err := gorm.G[Result](db).Raw("SELECT id, name, age FROM users WHERE id = ?", 3).Find(context.Background())

// Scan into a primitive type
age, err := gorm.G[int](db).Raw("SELECT SUM(age) FROM users WHERE role = ?", "admin").Find(context.Background())

// Scan into a slice
users, err := gorm.G[User](db).Raw("UPDATE users SET name = ? WHERE age = ? RETURNING id, name", "jinzhu", 20).Find(context.Background())
```

`Exec` with Raw SQL

```go
// Execute raw SQL
result := gorm.WithResult()
err := gorm.G[any](db, result).Exec(context.Background(), "DROP TABLE users")

// Execute with parameters
err = gorm.G[any](db).Exec(context.Background(), "UPDATE orders SET shipped_at = ? WHERE id IN ?", time.Now(), []int64{1, 2, 3})

// Exec with SQL Expression
err = gorm.G[any](db).Exec(context.Background(), "UPDATE users SET money = ? WHERE name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")
```

### Traditional API

Query Raw SQL with `Scan`

```go
type Result struct {
 ID int
 Name string
 Age int
}

var result Result
db.Raw("SELECT id, name, age FROM users WHERE id = ?", 3).Scan(&result)

db.Raw("SELECT id, name, age FROM users WHERE name = ?", "jinzhu").Scan(&result)

var age int
db.Raw("SELECT SUM(age) FROM users WHERE role = ?", "admin").Scan(&age)

var users []User
db.Raw("UPDATE users SET name = ? WHERE age = ? RETURNING id, name", "jinzhu", 20).Scan(&users)
```

`Exec` with Raw SQL

```go
db.Exec("DROP TABLE users")
db.Exec("UPDATE orders SET shipped_at = ? WHERE id IN ?", time.Now(), []int64{1, 2, 3})

// Exec with SQL Expression
db.Exec("UPDATE users SET money = ? WHERE name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")
```

{% note warn %}
**NOTE** GORM allows cache prepared statement to increase performance, checkout [Performance](performance.html) for details
{% endnote %}

## <span id="named_argument">Named Argument</span>

### Generics API

GORM supports named arguments with [`sql.NamedArg`](https://tip.golang.org/pkg/database/sql/#NamedArg), `map[string]interface{}{}` or struct, for example:

```go
users, err := gorm.G[User](db).Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(context.Background())
// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"

result3, err := gorm.G[User](db).Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu2"}).First(context.Background())
// SELECT * FROM `users` WHERE name1 = "jinzhu2" OR name2 = "jinzhu2" ORDER BY `users`.`id` LIMIT 1

// Named Argument with Raw SQL
users, err := gorm.G[User](db).Raw("SELECT * FROM users WHERE name1 = @name OR name2 = @name2 OR name3 = @name",
 sql.Named("name", "jinzhu1"), sql.Named("name2", "jinzhu2")).Find(context.Background())
// SELECT * FROM users WHERE name1 = "jinzhu1" OR name2 = "jinzhu2" OR name3 = "jinzhu1"

err := gorm.G[any](db).Exec(context.Background(), "UPDATE users SET name1 = @name, name2 = @name2, name3 = @name",
 sql.Named("name", "jinzhunew"), sql.Named("name2", "jinzhunew2"))
// UPDATE users SET name1 = "jinzhunew", name2 = "jinzhunew2", name3 = "jinzhunew"

users, err := gorm.G[User](db).Raw("SELECT * FROM users WHERE (name1 = @name AND name3 = @name) AND name2 = @name2",
 map[string]interface{}{"name": "jinzhu", "name2": "jinzhu2"}).Find(context.Background())
// SELECT * FROM users WHERE (name1 = "jinzhu" AND name3 = "jinzhu") AND name2 = "jinzhu2"

type NamedArgument struct {
	Name string
	Name2 string
}

users, err := gorm.G[User](db).Raw("SELECT * FROM users WHERE (name1 = @Name AND name3 = @Name) AND name2 = @Name2",
 NamedArgument{Name: "jinzhu", Name2: "jinzhu2"}).Find(context.Background())
// SELECT * FROM users WHERE (name1 = "jinzhu" AND name3 = "jinzhu") AND name2 = "jinzhu2"
```

### Traditional API

GORM supports named arguments with [`sql.NamedArg`](https://tip.golang.org/pkg/database/sql/#NamedArg), `map[string]interface{}{}` or struct, for example:

```go
db.Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(&user)
// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"

db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu2"}).First(&result3)
// SELECT * FROM `users` WHERE name1 = "jinzhu2" OR name2 = "jinzhu2" ORDER BY `users`.`id` LIMIT 1

// Named Argument with Raw SQL
db.Raw("SELECT * FROM users WHERE name1 = @name OR name2 = @name2 OR name3 = @name",
 sql.Named("name", "jinzhu1"), sql.Named("name2", "jinzhu2")).Find(&user)
// SELECT * FROM users WHERE name1 = "jinzhu1" OR name2 = "jinzhu2" OR name3 = "jinzhu1"

db.Exec("UPDATE users SET name1 = @name, name2 = @name2, name3 = @name",
 sql.Named("name", "jinzhunew"), sql.Named("name2", "jinzhunew2"))
// UPDATE users SET name1 = "jinzhunew", name2 = "jinzhunew2", name3 = "jinzhunew"

db.Raw("SELECT * FROM users WHERE (name1 = @name AND name3 = @name) AND name2 = @name2",
 map[string]interface{}{"name": "jinzhu", "name2": "jinzhu2"}).Find(&user)
// SELECT * FROM users WHERE (name1 = "jinzhu" AND name3 = "jinzhu") AND name2 = "jinzhu2"

type NamedArgument struct {
	Name string
	Name2 string
}

db.Raw("SELECT * FROM users WHERE (name1 = @Name AND name3 = @Name) AND name2 = @Name2",
 NamedArgument{Name: "jinzhu", Name2: "jinzhu2"}).Find(&user)
// SELECT * FROM users WHERE (name1 = "jinzhu" AND name3 = "jinzhu") AND name2 = "jinzhu2"
```

## DryRun Mode

Generate `SQL` and its arguments without executing, can be used to prepare or test generated SQL, Checkout [Session](session.html) for details

```go
stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
stmt.SQL.String() //=> SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
stmt.Vars //=> []interface{}{1}
```

## ToSQL

Returns generated `SQL` without executing.

During normal query execution, GORM uses `database/sql` argument placeholders to keep arguments separate from the SQL statement, which helps prevent SQL injection. However, the SQL generated by `ToSQL` does not provide the same safety guarantees and should only be used for debugging.

```go
sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
 return tx.Model(&User{}).Where("id = ?", 100).Limit(10).Order("age desc").Find(&[]User{})
})
sql //=> SELECT * FROM "users" WHERE id = 100 AND "users"."deleted_at" IS NULL ORDER BY age desc LIMIT 10
```

## `Row` & `Rows`

### Generics API

Get result as `*sql.Row`

```go
// Use GORM API build SQL
row := gorm.G[any](db).Table("users").Where("name = ?", "jinzhu").Select("name", "age").Row(context.Background())
row.Scan(&name, &age)

// Use Raw SQL
row := gorm.G[any](db).Raw("select name, age, email from users where name = ?", "jinzhu").Row(context.Background())
row.Scan(&name, &age, &email)
```

Get result as `*sql.Rows`

```go
// Use GORM API build SQL
rows, err := gorm.G[User](db).Where("name = ?", "jinzhu").Select("name, age, email").Rows(context.Background())
defer rows.Close()
for rows.Next() {
 rows.Scan(&name, &age, &email)

 // do something
}

// Raw SQL
rows, err := gorm.G[any](db).Raw("select name, age, email from users where name = ?", "jinzhu").Rows(context.Background())
defer rows.Close()
for rows.Next() {
 rows.Scan(&name, &age, &email)

 // do something
}
```

### Traditional API

Get result as `*sql.Row`

```go
// Use GORM API build SQL
row := db.Table("users").Where("name = ?", "jinzhu").Select("name", "age").Row()
row.Scan(&name, &age)

// Use Raw SQL
row := db.Raw("select name, age, email from users where name = ?", "jinzhu").Row()
row.Scan(&name, &age, &email)
```

Get result as `*sql.Rows`

```go
// Use GORM API build SQL
rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows()
defer rows.Close()
for rows.Next() {
 rows.Scan(&name, &age, &email)

 // do something
}

// Raw SQL
rows, err := db.Raw("select name, age, email from users where name = ?", "jinzhu").Rows()
defer rows.Close()
for rows.Next() {
 rows.Scan(&name, &age, &email)

 // do something
}
```

Checkout [FindInBatches](advanced_query.html) for how to query and process records in batch
Checkout [Group Conditions](advanced_query.html#group_conditions) for how to build complicated SQL Query

## Scan `*sql.Rows` into struct

Use `ScanRows` to scan a row into a struct, for example:

```go
rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows() // (*sql.Rows, error)
defer rows.Close()

var user User
for rows.Next() {
 // ScanRows scan a row into user
 db.ScanRows(rows, &user)

 // do something
}
```

## <span id="connection">Connection</span>

Run mutliple SQL in same db tcp connection (not in a transaction)

```go
db.Connection(func(tx *gorm.DB) error {
 tx.Exec("SET my.role = ?", "admin")

 tx.First(&User{})
})
```

## Advanced

### <span id="clauses">Clauses</span>

GORM uses SQL builder generates SQL internally, for each operation, GORM creates a `*gorm.Statement` object, all GORM APIs add/change `Clause` for the `Statement`, at last, GORM generated SQL based on those clauses

For example, when querying with `First`, it adds the following clauses to the `Statement`

```go
var limit = 1
clause.Select{Columns: []clause.Column{{Name: "*"}}}
clause.From{Tables: []clause.Table{{Name: clause.CurrentTable}}}
clause.Limit{Limit: &limit}
clause.OrderBy{Columns: []clause.OrderByColumn{
 {
 Column: clause.Column{
 Table: clause.CurrentTable,
 Name: clause.PrimaryKey,
 },
 },
}}
```

Then GORM build finally querying SQL in the `Query` callbacks like:

```go
Statement.Build("SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT", "FOR")
```

Which generate SQL:

```sql
SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1
```

You can define your own `Clause` and use it with GORM, it needs to implements [Interface](https://pkg.go.dev/gorm.io/gorm/clause?tab=doc#Interface)

Check out [examples](https://github.com/go-gorm/gorm/tree/master/clause) for reference

### Clause Builder

For different databases, Clauses may generate different SQL, for example:

```go
db.Offset(10).Limit(5).Find(&users)
// Generated for SQL Server
// SELECT * FROM "users" OFFSET 10 ROW FETCH NEXT 5 ROWS ONLY
// Generated for MySQL
// SELECT * FROM `users` LIMIT 5 OFFSET 10
```

Which is supported because GORM allows database driver register Clause Builder to replace the default one, take the [Limit](https://github.com/go-gorm/sqlserver/blob/512546241200023819d2e7f8f2f91d7fb3a52e42/sqlserver.go#L45) as example

### Clause Options

GORM defined [Many Clauses](https://github.com/go-gorm/gorm/tree/master/clause), and some clauses provide advanced options can be used for your application

Although most of them are rarely used, if you find GORM public API can't match your requirements, may be good to check them out, for example:

```go
db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&user)
// INSERT IGNORE INTO users (name,age...) VALUES ("jinzhu",18...);
```

### StatementModifier

GORM provides interface [StatementModifier](https://pkg.go.dev/gorm.io/gorm?tab=doc#StatementModifier) allows you modify statement to match your requirements, take [Hints](hints.html) as example

```go
import "gorm.io/hints"

db.Clauses(hints.New("hint")).Find(&User{})
// SELECT * /*+ hint */ FROM `users`
```

---
title: The Generics Way to Use GORM
layout: page
---

GORM has officially introduced support for **Go Generics** in its latest version (>= `v1.30.0`). This addition significantly enhances usability and type safety while reducing issues such as SQL pollution caused by reusing `gorm.DB` instances. Additionally, we've improved the behaviors of `Joins` and `Preload` and incorporated transaction timeout handling to prevent connection pool leaks.

This update introduces generic APIs in a carefully designed way that maintains full backward compatibility with existing APIs. You can freely mix traditional and generic APIs in your projects—just use generics for new code without worrying about compatibility with existing logic or GORM plugins (such as encryption/decryption, sharding, read/write splitting, tracing, etc.).

To prevent misuse, we have intentionally removed certain APIs in the generics version that are prone to ambiguity or concurrency issues, such as `FirstOrCreate` and `Save`. At the same time, we are designing a brand new `gorm` CLI tool, which will offer stronger code generation capabilities, enhanced type safety, and lint support in the future — further reducing the risk of incorrect usage.

We strongly recommend using the new generics-based API in new projects or during refactoring efforts to enjoy a better development experience, improved type guarantees, and a more maintainable codebase.

## Generic APIs

GORM's generic APIs closely mirror the functionality of the original ones. Here are some common operations using the new generics APIs:

```go
ctx := context.Background()

// Create records
gorm.G[User](db).Create(ctx, &User{Name: "Alice"})
gorm.G[User](db).CreateInBatches(ctx, users, 10)

// Query records
user, err := gorm.G[User](db).Where("name = ?", "Jinzhu").First(ctx)
users, err := gorm.G[User](db).Where("age <= ?", 18).Find(ctx)

// Update records
gorm.G[User](db).Where("id = ?", u.ID).Update(ctx, "age", 18)
gorm.G[User](db).Where("id = ?", u.ID).Updates(ctx, User{Name: "Jinzhu", Age: 18})

// Delete records
gorm.G[User](db).Where("id = ?", u.ID).Delete(ctx)
```

The generics APIs fully support GORM’s advanced features by accepting optional parameters, such as clause configurations or plugin-based options (e.g., hints, resolvers), enabling powerful and flexible behaviors.

```go
// OnConflict: Handle conflict during insert
err := gorm.G[Language](DB, clause.OnConflict{DoNothing: true}).Create(ctx, &lang)
err := gorm.G[Language](DB, clause.OnConflict{
 Columns: []clause.Column{{Name: "id"}},
 DoUpdates: clause.Assignments(map[string]interface{}{"count": gorm.Expr("GREATEST(count, VALUES(count))")}),
}).Create(ctx, &lang)

// Execution hints
err := gorm.G[User](DB,
 hints.New("MAX_EXECUTION_TIME(100)"),
 hints.New("USE_INDEX(t1, idx1)"),
).Find(ctx)
// SELECT /*+ MAX_EXECUTION_TIME(100) USE_INDEX(t1, idx1) */ * FROM `users`

// Read from master in read/write splitting mode
err := gorm.G[User](DB, dbresolver.Write).Find(ctx)

// Retrieve raw result metadata
result := gorm.WithResult()
err := gorm.G[User](DB, result).CreateInBatches(ctx, &users, 2)
// result.RowsAffected
// result.Result.LastInsertId()
```

## Joins / Preload Enhancements

The new GORM generics interface brings enhanced support for association queries (`Joins`) and eager loading (`Preload`), offering more flexible association methods, more expressive query capabilities, and a significantly simplified approach to building complex queries.

* **Joins**: Easily specify different join types (e.g., `InnerJoin`, `LeftJoin`) and customize join conditions based on associations, making complex cross-table queries clearer and more intuitive.

```go
// Load only users who have a company
users, err := gorm.G[User](db).Joins(clause.Has("Company"), nil).Find(ctx)

// Use Left Join with custom filter on joined table
user, err = gorm.G[User](db).Joins(clause.LeftJoin.Association("Company"), func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
 db.Where(map[string]any{"name": company.Name})
 return nil
}).Where(map[string]any{"name": user.Name}).First(ctx)

// Join using a subquery
users, err = gorm.G[User](db).Joins(clause.LeftJoin.AssociationFrom("Company", gorm.G[Company](DB).Select("Name")).As("t"),
 func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
 db.Where("?.name = ?", joinTable, u.Company.Name)
 return nil
 },
).Find(ctx)
```

* **Preload**: Simplifies conditions for eager loading and introduces the `LimitPerRecord` option, which allows limiting the number of related records loaded per primary record when eager loading collections.

```go
// A basic Preload example
users, err := gorm.G[User](db).Preload("Friends", func(db gorm.PreloadBuilder) error {
 db.Where("age > ?", 14)
 return nil
}).Where("age > ?", 18).Find(ctx)

// Preload nested associations
users, err := gorm.G[User](db).Preload("Friends.Pets", nil).Where("age > ?", 18).Find(ctx)

// Preload with sort and per-record limit
users, err = gorm.G[User](db).Preload("Friends", func(db gorm.PreloadBuilder) error {
 db.Select("id", "name").Order("age desc")
 return nil
}).Preload("Friends.Pets", func(db gorm.PreloadBuilder) error {
 db.LimitPerRecord(2)
 return nil
}).Find(ctx)
```

## Complex Raw SQL

The generics interface continues to support `Raw` SQL execution for complex or edge-case scenarios:

```go
users, err := gorm.G[User](DB).Raw("SELECT name FROM users WHERE id = ?", user.ID).Find(ctx)
```

However, we **strongly recommend** using our new **code generation tool** to achieve type-safe, maintainable, and secure raw queries—reducing risks like syntax errors or SQL injection.

### Code Generator Workflow

* **1. Install the CLI tool:**

```bash
go install gorm.io/cli/gorm@latest
```

* **2. Define query interfaces:**

Simply define your query interface using Go’s `interface` syntax, embedding SQL templates as comments:

```go
type Query[T any] interface {
	// GetByID queries data by ID and returns it as a struct.
	//
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id int) (T, error)

	// SELECT * FROM @@table WHERE @@column=@value
	FilterWithColumn(column string, value string) (T, error)

	// SELECT * FROM users
	// {{if user.ID > 0}}
	// WHERE id=@user.ID
	// {{else if user.Name != ""}}
	// WHERE username=@user.Name
	// {{end}}
	QueryWith(user models.User) (T, error)

	// UPDATE @@table
	// {{set}}
	// {{if user.Name != ""}} username=@user.Name, {{end}}
	// {{if user.Age > 0}} age=@user.Age, {{end}}
	// {{if user.Age >= 18}} is_adult=1 {{else}} is_adult=0 {{end}}
	// {{end}}
	// WHERE id=@id
	Update(user models.User, id int) error

	// SELECT * FROM @@table
	// {{where}}
	// {{for _, user := range users}}
	// {{if user.Name != "" && user.Age > 0}}
	// (username = @user.Name AND age=@user.Age AND role LIKE concat("%",@user.Role,"%")) OR
	// {{end}}
	// {{end}}
	// {{end}}
	Filter(users []models.User) ([]T, error)

	// where("name=@name AND age=@age")
	FilterByNameAndAge(name string, age int)

	// SELECT * FROM @@table
	// {{where}}
	// {{if !start.IsZero()}}
	// created_time > @start
	// {{end}}
	// {{if !end.IsZero()}}
	// AND created_time < @end
	// {{end}}
	// {{end}}
	FilterWithTime(start, end time.Time) ([]T, error)
}
```

* **3. Run the generator:**

```bash
gorm gen -i ./examples/example.go -o query
```

* **4. Use the generated API:**

```go
import "your_project/query"

company, err := query.Query[Company](db).GetByID(ctx, 10)
// SELECT * FROM `companies` WHERE id=10
user, err := query.Query[User](db).GetByID(ctx, 10)
// SELECT * FROM `users` WHERE id=10

// Combine with other Generic APIs
err := query.Query[User](db).FilterByNameAndAge("jinzhu", 18).Delete(ctx)
// DELETE FROM `users` WHERE name='jinzhu' AND age=18

users, err := query.Query[User](db).FilterByNameAndAge("jinzhu", 18).Find(ctx)
// SELECT * FROM `users` WHERE name='jinzhu' AND age=18
```

## Summary

This release marks a significant step forward for GORM in both generics support and the brand-new `gorm` command-line tool. These features have been in the planning stage for quite some time, and we’re excited to finally bring an initial implementation to the community.

In the coming updates, we’ll continue refining the generics API, enhancing the CLI tool, and updating the official [https://gorm.io](https://gorm.io) documentation accordingly—aiming to provide a clearer, more efficient developer experience.

We deeply appreciate the support from all GORM users and sponsors over the years. GORM’s growth over the past 12 years simply wouldn’t have been possible without you ❤️

---
title: Transactions
layout: page
---

## Disable Default Transaction

GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency, you can disable it during initialization if it is not required, you will gain about 30%+ performance improvement after that

```go
// Globally disable
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 SkipDefaultTransaction: true,
})

// Continuous session mode
tx := db.Session(&Session{SkipDefaultTransaction: true})
tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)
```

## Transaction

To perform a set of operations within a transaction, the general flow is as below.

### Generics API

```go
ctx := context.Background()

// Basic transaction
err := db.Transaction(func(tx *gorm.DB) error {
 // Use Generics API inside the transaction
 if err := gorm.G[Animal](tx).Create(ctx, &Animal{Name: "Giraffe"}); err != nil {
 // return any error will rollback
 return err
 }

 if err := gorm.G[Animal](tx).Create(ctx, &Animal{Name: "Lion"}); err != nil {
 return err
 }

 // return nil will commit the whole transaction
 return nil
})
```

### Traditional API

```go
db.Transaction(func(tx *gorm.DB) error {
 // do some database operations in the transaction (use 'tx' from this point, not 'db')
 if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
 // return any error will rollback
 return err
 }

 if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
 return err
 }

 // return nil will commit the whole transaction
 return nil
})
```

### Nested Transactions

GORM supports nested transactions, you can rollback a subset of operations performed within the scope of a larger transaction, for example:

```go
db.Transaction(func(tx *gorm.DB) error {
 tx.Create(&user1)

 tx.Transaction(func(tx2 *gorm.DB) error {
 tx2.Create(&user2)
 return errors.New("rollback user2") // Rollback user2
 })

 tx.Transaction(func(tx3 *gorm.DB) error {
 tx3.Create(&user3)
 return nil
 })

 return nil
})

// Commit user1, user3
```

## Control the transaction manually

Gorm supports calling transaction control functions (commit / rollback) directly, for example:

```go
// begin a transaction
tx := db.Begin()

// do some database operations in the transaction (use 'tx' from this point, not 'db')
tx.Create(...)

// ...

// rollback the transaction in case of error
tx.Rollback()

// Or commit the transaction
tx.Commit()
```

### A Specific Example

```go
func CreateAnimals(db *gorm.DB) error {
 // Note the use of tx as the database handle once you are within a transaction
 tx := db.Begin()
 defer func() {
 if r := recover(); r != nil {
 tx.Rollback()
 }
 }()

 if err := tx.Error; err != nil {
 return err
 }

 if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
 tx.Rollback()
 return err
 }

 if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
 tx.Rollback()
 return err
 }

 return tx.Commit().Error
}
```

## SavePoint, RollbackTo

GORM provides `SavePoint`, `RollbackTo` to save points and roll back to a savepoint, for example:

```go
tx := db.Begin()
tx.Create(&user1)

tx.SavePoint("sp1")
tx.Create(&user2)
tx.RollbackTo("sp1") // Rollback user2

tx.Commit() // Commit user1
```

---
title: Update
layout: page
---

## Save All Fields

### Traditional API

`Save` will save all fields when performing the Updating SQL

```go
db.First(&user)

user.Name = "jinzhu 2"
user.Age = 100
db.Save(&user)
// UPDATE users SET name='jinzhu 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=111;
```

`Save` is an upsert function:
- If the value contains no primary key, it performs `Create`
- If the value has a primary key, it first executes **Update** (all fields, by `Select(*)`).
- If `rows affected = 0` after **Update**, it automatically falls back to `Create`.

> 💡 **Note**: `Save` guarantees either an update or insert will occur.
> To prevent unintended creation when no rows match, use [ `Select(*).Updates()` ](update.html#Update-Selected-Fields).

```go
db.Save(&User{Name: "jinzhu", Age: 100})
// INSERT INTO `users` (`name`,`age`,`birthday`,`update_at`) VALUES ("jinzhu",100,"0000-00-00 00:00:00","0000-00-00 00:00:00")

db.Save(&User{ID: 1, Name: "jinzhu", Age: 100})
// UPDATE `users` SET `name`="jinzhu",`age`=100,`birthday`="0000-00-00 00:00:00",`update_at`="0000-00-00 00:00:00" WHERE `id` = 1
```

{% note warn %}
**NOTE** Don't use `Save` with `Model`, it's an **Undefined Behavior**.
{% endnote %}

{% note warn %}
**NOTE** The `Save` method is intentionally removed from the Generics API to prevent ambiguity and concurrency issues. Please use `Create` or `Updates` methods instead.
{% endnote %}

## Update single column

When updating a single column with `Update`, it needs to have any conditions or it will raise error `ErrMissingWhereClause`, checkout [Block Global Updates](#block_global_updates) for details.

### Generics API

```go
ctx := context.Background()

// Update with conditions
err := gorm.G[User](db).Where("active = ?", true).Update(ctx, "name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;

// Update with ID condition
err := gorm.G[User](db).Where("id = ?", 111).Update(ctx, "name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

// Update with multiple conditions
err := gorm.G[User](db).Where("id = ? AND active = ?", 111, true).Update(ctx, "name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;
```

### Traditional API

When using the `Model` method and its value has a primary value, the primary key will be used to build the condition, for example:

```go
// Update with conditions
db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;

// User's ID is `111`:
db.Model(&user).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

// Update with conditions and model value
db.Model(&user).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;
```

## Updates multiple columns

`Updates` supports updating with `struct` or `map[string]interface{}`, when updating with `struct` it will only update non-zero fields by default

### Generics API

```go
ctx := context.Background()

// Update attributes with `struct`, will only update non-zero fields
rows, err := gorm.G[User](db).Where("id = ?", 111).Updates(ctx, User{Name: "hello", Age: 18, Active: false})
// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

// Update attributes with `map`
rows, err := gorm.G[map[string]interface{}](db).Table("users").Where("id = ?", 111).Updates(ctx, map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
```

### Traditional API

```go
// Update attributes with `struct`, will only update non-zero fields
db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

// Update attributes with `map`
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
```

{% note warn %}
**NOTE** When updating with struct, GORM will only update non-zero fields. You might want to use `map` to update attributes or use `Select` to specify fields to update
{% endnote %}

## Update Selected Fields

If you want to update selected fields or ignore some fields when updating, you can use `Select`, `Omit`

### Generics API

```go
ctx := context.Background()

// Select with Map
rows, err := gorm.G[User](db).Where("id = ?", 111).Select("name").Updates(ctx, map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello' WHERE id=111;

rows, err := gorm.G[User](db).Where("id = ?", 111).Omit("name").Updates(ctx, map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

// Select with Struct (select zero value fields)
rows, err := gorm.G[User](db).Where("id = ?", 111).Select("Name", "Age").Updates(ctx, User{Name: "new_name", Age: 0})
// UPDATE users SET name='new_name', age=0 WHERE id=111;

// Select all fields (select all fields include zero value fields)
rows, err := gorm.G[User](db).Where("id = ?", 111).Select("*").Updates(ctx, User{Name: "jinzhu", Role: "admin", Age: 0})

// Select all fields but omit Role (select all fields include zero value fields)
rows, err := gorm.G[User](db).Where("id = ?", 111).Select("*").Omit("Role").Updates(ctx, User{Name: "jinzhu", Role: "admin", Age: 0})
```

### Traditional API

```go
// Select with Map
// User's ID is `111`:
db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello' WHERE id=111;

db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

// Select with Struct (select zero value fields)
db.Model(&user).Select("Name", "Age").Updates(User{Name: "new_name", Age: 0})
// UPDATE users SET name='new_name', age=0 WHERE id=111;

// Select all fields (select all fields include zero value fields)
db.Model(&user).Select("*").Updates(User{Name: "jinzhu", Role: "admin", Age: 0})

// Select all fields but omit Role (select all fields include zero value fields)
db.Model(&user).Select("*").Omit("Role").Updates(User{Name: "jinzhu", Role: "admin", Age: 0})
```

## Update Hooks

GORM allows the hooks `BeforeSave`, `BeforeUpdate`, `AfterSave`, `AfterUpdate`. Those methods will be called when updating a record, refer [Hooks](hooks.html) for details

```go
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Role == "admin" {
 return errors.New("admin user not allowed to update")
	}
	return
}
```

## Batch Updates

If we haven't specified a record having a primary key value with `Model`, GORM will perform a batch update

```go
// Update with struct
db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
// UPDATE users SET name='hello', age=18 WHERE role = 'admin';

// Update with map
db.Table("users").Where("id IN ?", []int{10, 11}).Updates(map[string]interface{}{"name": "hello", "age": 18})
// UPDATE users SET name='hello', age=18 WHERE id IN (10, 11);
```

### <span id="block_global_updates">Block Global Updates</span>

If you perform a batch update without any conditions, GORM WON'T run it and will return `ErrMissingWhereClause` error by default

You have to use some conditions or use raw SQL or enable the `AllowGlobalUpdate` mode, for example:

```go
db.Model(&User{}).Update("name", "jinzhu").Error // gorm.ErrMissingWhereClause

db.Model(&User{}).Where("1 = 1").Update("name", "jinzhu")
// UPDATE users SET `name` = "jinzhu" WHERE 1=1

db.Exec("UPDATE users SET name = ?", "jinzhu")
// UPDATE users SET name = "jinzhu"

db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&User{}).Update("name", "jinzhu")
// UPDATE users SET `name` = "jinzhu"
```

### Updated Records Count

Get the number of rows affected by a update

```go
// Get updated records count with `RowsAffected`
result := db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
// UPDATE users SET name='hello', age=18 WHERE role = 'admin';

result.RowsAffected // returns updated records count
result.Error // returns updating error
```

## Advanced

### <span id="update_from_sql_expr">Update with SQL Expression</span>

GORM allows updating a column with a SQL expression, e.g:

```go
// product's ID is `3`
db.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))
// UPDATE "products" SET "price" = price * 2 + 100, "updated_at" = '2013-11-17 21:34:10' WHERE "id" = 3;

db.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})
// UPDATE "products" SET "price" = price * 2 + 100, "updated_at" = '2013-11-17 21:34:10' WHERE "id" = 3;

db.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = 3;

db.Model(&product).Where("quantity > 1").UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = 3 AND quantity > 1;
```

And GORM also allows updating with SQL Expression/Context Valuer with [Customized Data Types](data_types.html#gorm_valuer_interface), e.g:

```go
// Create from customized data type
type Location struct {
	X, Y int
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
 return clause.Expr{
 SQL: "ST_PointFromText(?)",
 Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y)},
 }
}

db.Model(&User{ID: 1}).Updates(User{
 Name: "jinzhu",
 Location: Location{X: 100, Y: 100},
})
// UPDATE `user_with_points` SET `name`="jinzhu",`location`=ST_PointFromText("POINT(100 100)") WHERE `id` = 1
```

### Update from SubQuery

Update a table by using SubQuery

```go
db.Model(&user).Update("company_name", db.Model(&Company{}).Select("name").Where("companies.id = users.company_id"))
// UPDATE "users" SET "company_name" = (SELECT name FROM companies WHERE companies.id = users.company_id);

db.Table("users as u").Where("name = ?", "jinzhu").Update("company_name", db.Table("companies as c").Select("name").Where("c.id = u.company_id"))

db.Table("users as u").Where("name = ?", "jinzhu").Updates(map[string]interface{}{"company_name": db.Table("companies as c").Select("name").Where("c.id = u.company_id")})
```

### Without Hooks/Time Tracking

If you want to skip `Hooks` methods and don't track the update time when updating, you can use `UpdateColumn`, `UpdateColumns`, it works like `Update`, `Updates`

```go
// Update single column
db.Model(&user).UpdateColumn("name", "hello")
// UPDATE users SET name='hello' WHERE id = 111;

// Update multiple columns
db.Model(&user).UpdateColumns(User{Name: "hello", Age: 18})
// UPDATE users SET name='hello', age=18 WHERE id = 111;

// Update selected columns
db.Model(&user).Select("name", "age").UpdateColumns(User{Name: "hello", Age: 0})
// UPDATE users SET name='hello', age=0 WHERE id = 111;
```

### Returning Data From Modified Rows

Returning changed data only works for databases which support Returning, for example:

```go
// return all columns
var users []User
db.Model(&users).Clauses(clause.Returning{}).Where("role = ?", "admin").Update("salary", gorm.Expr("salary * ?", 2))
// UPDATE `users` SET `salary`=salary * 2,`updated_at`="2021-10-28 17:37:23.19" WHERE role = "admin" RETURNING *
// users => []User{{ID: 1, Name: "jinzhu", Role: "admin", Salary: 100}, {ID: 2, Name: "jinzhu.2", Role: "admin", Salary: 1000}}

// return specified columns
db.Model(&users).Clauses(clause.Returning{Columns: []clause.Column{{Name: "name"}, {Name: "salary"}}}).Where("role = ?", "admin").Update("salary", gorm.Expr("salary * ?", 2))
// UPDATE `users` SET `salary`=salary * 2,`updated_at`="2021-10-28 17:37:23.19" WHERE role = "admin" RETURNING `name`, `salary`
// users => []User{{ID: 0, Name: "jinzhu", Role: "", Salary: 100}, {ID: 0, Name: "jinzhu.2", Role: "", Salary: 1000}}
```

### Check Field has changed?

GORM provides the `Changed` method which could be used in **Before Update Hooks**, it will return whether the field has changed or not.

The `Changed` method only works with methods `Update`, `Updates`, and it only checks if the updating value from `Update` / `Updates` equals the model value. It will return true if it is changed and not omitted

```go
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
 // if Role changed
	if tx.Statement.Changed("Role") {
 return errors.New("role not allowed to change")
	}

 if tx.Statement.Changed("Name", "Admin") { // if Name or Role changed
 tx.Statement.SetColumn("Age", 18)
 }

 // if any fields changed
	if tx.Statement.Changed() {
 tx.Statement.SetColumn("RefreshedAt", time.Now())
	}
	return nil
}

db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(map[string]interface{"name": "jinzhu2"})
// Changed("Name") => true
db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(map[string]interface{"name": "jinzhu"})
// Changed("Name") => false, `Name` not changed
db.Model(&User{ID: 1, Name: "jinzhu"}).Select("Admin").Updates(map[string]interface{
 "name": "jinzhu2", "admin": false,
})
// Changed("Name") => false, `Name` not selected to update

db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(User{Name: "jinzhu2"})
// Changed("Name") => true
db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(User{Name: "jinzhu"})
// Changed("Name") => false, `Name` not changed
db.Model(&User{ID: 1, Name: "jinzhu"}).Select("Admin").Updates(User{Name: "jinzhu2"})
// Changed("Name") => false, `Name` not selected to update
```

### Change Updating Values

To change updating values in Before Hooks, you should use `SetColumn` unless it is a full update with `Save`, for example:

```go
func (user *User) BeforeSave(tx *gorm.DB) (err error) {
 if pw, err := bcrypt.GenerateFromPassword(user.Password, 0); err == nil {
 tx.Statement.SetColumn("EncryptedPassword", pw)
 }

 if tx.Statement.Changed("Code") {
 user.Age += 20
 tx.Statement.SetColumn("Age", user.Age)
 }
}

db.Model(&user).Update("Name", "jinzhu")
```

---
title: GORM 2.0 Release Note
layout: page
---

GORM 2.0 is a rewrite from scratch, it introduces some incompatible-API change and many improvements

**Highlights**

* Performance Improvements
* Modularity
* Context, Batch Insert, Prepared Statement Mode, DryRun Mode, Join Preload, Find To Map, Create From Map, FindInBatches supports
* Nested Transaction/SavePoint/RollbackTo SavePoint supports
* SQL Builder, Named Argument, Group Conditions, Upsert, Locking, Optimizer/Index/Comment Hints supports, SubQuery improvements, CRUD with SQL Expr and Context Valuer
* Full self-reference relationships support, Join Table improvements, Association Mode for batch data
* Multiple fields allowed to track create/update time, UNIX (milli/nano) seconds supports
* Field permissions support: read-only, write-only, create-only, update-only, ignored
* New plugin system, provides official plugins for multiple databases, read/write splitting, prometheus integrations...
* New Hooks API: unified interface with plugins
* New Migrator: allows to create database foreign keys for relationships, smarter AutoMigrate, constraints/checker support, enhanced index support
* New Logger: context support, improved extensibility
* Unified Naming strategy: table name, field name, join table name, foreign key, checker, index name rules
* Better customized data type support (e.g: JSON)

## How To Upgrade

* GORM's developments moved to [github.com/go-gorm](https://github.com/go-gorm), and its import path changed to `gorm.io/gorm`, for previous projects, you can keep using `github.com/jinzhu/gorm` [GORM V1 Document](https://v1.gorm.io/)
* Database drivers have been split into separate projects, e.g: [github.com/go-gorm/sqlite](https://github.com/go-gorm/sqlite), and its import path also changed to `gorm.io/driver/sqlite`

### Install

```go
go get gorm.io/gorm
// **NOTE** GORM `v2.0.0` released with git tag `v1.20.0`
```

### Quick Start

```go
import (
 "gorm.io/gorm"
 "gorm.io/driver/sqlite"
)

func init() {
 db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

 // Most CRUD API kept compatibility
 db.AutoMigrate(&Product{})
 db.Create(&user)
 db.First(&user, 1)
 db.Model(&user).Update("Age", 18)
 db.Model(&user).Omit("Role").Updates(map[string]interface{}{"Name": "jinzhu", "Role": "admin"})
 db.Delete(&user)
}
```

## Major Features

The release note only cover major changes introduced in GORM V2 as a quick reference list

#### Context Support

* Database operations support `context.Context` with the `WithContext` method
* Logger also accepts context for tracing

```go
db.WithContext(ctx).Find(&users)
```

#### Batch Insert

To efficiently insert large number of records, pass a slice to the `Create` method. GORM will generate a single SQL statement to insert all the data and backfill primary key values, hook methods will be invoked too.

```go
var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
db.Create(&users)

for _, user := range users {
 user.ID // 1,2,3
}
```

You can specify batch size when creating with `CreateInBatches`, e.g:

```go
var users = []User{{Name: "jinzhu_1"}, ...., {Name: "jinzhu_10000"}}

// batch size 100
db.CreateInBatches(users, 100)
```

#### Prepared Statement Mode

Prepared Statement Mode creates prepared stmt and caches them to speed up future calls

```go
// globally mode, all operations will create prepared stmt and cache to speed up
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{PrepareStmt: true})

// session mode, create prepares stmt and speed up current session operations
tx := db.Session(&Session{PrepareStmt: true})
tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)
```

#### DryRun Mode

Generates SQL without executing, can be used to check or test generated SQL

```go
stmt := db.Session(&Session{DryRun: true}).Find(&user, 1).Statement
stmt.SQL.String() //=> SELECT * FROM `users` WHERE `id` = $1 // PostgreSQL
stmt.SQL.String() //=> SELECT * FROM `users` WHERE `id` = ? // MySQL
stmt.Vars //=> []interface{}{1}
```

#### Join Preload

Preload associations using INNER JOIN, and will handle null data to avoid failing to scan

```go
db.Joins("Company").Joins("Manager").Joins("Account").Find(&users, "users.id IN ?", []int{1,2})
```

#### Find To Map

Scan result to `map[string]interface{}` or `[]map[string]interface{}`

```go
var result map[string]interface{}
db.Model(&User{}).First(&result, "id = ?", 1)
```

#### Create From Map

Create from map `map[string]interface{}` or `[]map[string]interface{}`

```go
db.Model(&User{}).Create(map[string]interface{}{"Name": "jinzhu", "Age": 18})

datas := []map[string]interface{}{
 {"Name": "jinzhu_1", "Age": 19},
 {"name": "jinzhu_2", "Age": 20},
}

db.Model(&User{}).Create(datas)
```

#### FindInBatches

Query and process records in batch

```go
result := db.Where("age>?", 13).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
 // batch processing
 return nil
})
```

#### Nested Transaction

```go
db.Transaction(func(tx *gorm.DB) error {
 tx.Create(&user1)

 tx.Transaction(func(tx2 *gorm.DB) error {
 tx.Create(&user2)
 return errors.New("rollback user2") // rollback user2
 })

 tx.Transaction(func(tx2 *gorm.DB) error {
 tx.Create(&user3)
 return nil
 })

 return nil // commit user1 and user3
})
```

#### SavePoint, RollbackTo

```go
tx := db.Begin()
tx.Create(&user1)

tx.SavePoint("sp1")
tx.Create(&user2)
tx.RollbackTo("sp1") // rollback user2

tx.Commit() // commit user1
```

#### Named Argument

GORM supports use `sql.NamedArg`, `map[string]interface{}` as named arguments

```go
db.Where("name1 = @name OR name2 = @name", sql.Named("name", "jinzhu")).Find(&user)
// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu"

db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu2"}).First(&result3)
// SELECT * FROM `users` WHERE name1 = "jinzhu2" OR name2 = "jinzhu2" ORDER BY `users`.`id` LIMIT 1

db.Raw(
 "SELECT * FROM users WHERE name1 = @name OR name2 = @name2 OR name3 = @name",
 sql.Named("name", "jinzhu1"), sql.Named("name2", "jinzhu2"),
).Find(&user)
// SELECT * FROM users WHERE name1 = "jinzhu1" OR name2 = "jinzhu2" OR name3 = "jinzhu1"

db.Exec(
 "UPDATE users SET name1 = @name, name2 = @name2, name3 = @name",
 map[string]interface{}{"name": "jinzhu", "name2": "jinzhu2"},
)
// UPDATE users SET name1 = "jinzhu", name2 = "jinzhu2", name3 = "jinzhu"
```

#### Group Conditions

```go
db.Where(
 db.Where("pizza = ?", "pepperoni").Where(db.Where("size = ?", "small").Or("size = ?", "medium")),
).Or(
 db.Where("pizza = ?", "hawaiian").Where("size = ?", "xlarge"),
).Find(&pizzas)

// SELECT * FROM pizzas WHERE (pizza = 'pepperoni' AND (size = 'small' OR size = 'medium')) OR (pizza = 'hawaiian' AND size = 'xlarge')
```

#### SubQuery

```go
// Where SubQuery
db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)

// From SubQuery
db.Table("(?) as u", db.Model(&User{}).Select("name", "age")).Where("age = ?", 18}).Find(&User{})
// SELECT * FROM (SELECT `name`,`age` FROM `users`) as u WHERE age = 18

// Update SubQuery
db.Model(&user).Update(
 "price", db.Model(&Company{}).Select("name").Where("companies.id = users.company_id"),
)
```

#### Upsert

`clause.OnConflict` provides compatible Upsert support for different databases (SQLite, MySQL, PostgreSQL, SQL Server)

```go
import "gorm.io/gorm/clause"

db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users)

db.Clauses(clause.OnConflict{
 Columns: []clause.Column{{Name: "id"}},
 DoUpdates: clause.Assignments(map[string]interface{}{"name": "jinzhu", "age": 18}),
}).Create(&users)
// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET ***; SQL Server
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE name="jinzhu", age=18; MySQL

db.Clauses(clause.OnConflict{
 Columns: []clause.Column{{Name: "id"}},
 DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
}).Create(&users)
// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET "name"="excluded"."name"; SQL Server
// INSERT INTO "users" *** ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name", "age"="excluded"."age"; PostgreSQL
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `name`=VALUES(name),`age=VALUES(age); MySQL
```

#### Locking

```go
db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
// SELECT * FROM `users` FOR UPDATE

db.Clauses(clause.Locking{
 Strength: "SHARE",
 Table: clause.Table{Name: clause.CurrentTable},
}).Find(&users)
// SELECT * FROM `users` FOR SHARE OF `users`
```

#### Optimizer/Index/Comment Hints

```go
import "gorm.io/hints"

// Optimizer Hints
db.Clauses(hints.New("hint")).Find(&User{})
// SELECT * /*+ hint */ FROM `users`

// Index Hints
db.Clauses(hints.UseIndex("idx_user_name")).Find(&User{})
// SELECT * FROM `users` USE INDEX (`idx_user_name`)

// Comment Hints
db.Clauses(hints.Comment("select", "master")).Find(&User{})
// SELECT /*master*/ * FROM `users`;
```

Check out [Hints](hints.html) for details

#### CRUD From SQL Expr/Context Valuer

```go
type Location struct {
	X, Y int
}

func (loc Location) GormDataType() string {
 return "geometry"
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
 return clause.Expr{
 SQL: "ST_PointFromText(?)",
 Vars: []interface{}{fmt.Sprintf("POINT(%d %d)", loc.X, loc.Y)},
 }
}

db.Create(&User{
 Name: "jinzhu",
 Location: Location{X: 100, Y: 100},
})
// INSERT INTO `users` (`name`,`point`) VALUES ("jinzhu",ST_PointFromText("POINT(100 100)"))

db.Model(&User{ID: 1}).Updates(User{
 Name: "jinzhu",
 Point: Point{X: 100, Y: 100},
})
// UPDATE `user_with_points` SET `name`="jinzhu",`point`=ST_PointFromText("POINT(100 100)") WHERE `id` = 1
```

Check out [Customize Data Types](data_types.html#gorm_valuer_interface) for details

#### Field permissions

Field permissions support, permission levels: read-only, write-only, create-only, update-only, ignored

```go
type User struct {
 Name string `gorm:"<-:create"` // allow read and create
 Name string `gorm:"<-:update"` // allow read and update
 Name string `gorm:"<-"` // allow read and write (create and update)
 Name string `gorm:"->:false;<-:create"` // createonly
 Name string `gorm:"->"` // readonly
 Name string `gorm:"-"` // ignored
}
```

#### Track creating/updating time/unix (milli/nano) seconds for multiple fields

```go
type User struct {
 CreatedAt time.Time // Set to current time if it is zero on creating
 UpdatedAt int // Set to current unix seconds on updaing or if it is zero on creating
 Updated int64 `gorm:"autoUpdateTime:nano"` // Use unix Nano seconds as updating time
 Updated2 int64 `gorm:"autoUpdateTime:milli"` // Use unix Milli seconds as updating time
 Created int64 `gorm:"autoCreateTime"` // Use unix seconds as creating time
}
```

#### Multiple Databases, Read/Write Splitting

GORM provides multiple databases, read/write splitting support with plugin `DB Resolver`, which also supports auto-switching database/table based on current struct/table, and multiple sources、replicas supports with customized load-balancing logic

Check out [Database Resolver](dbresolver.html) for details

#### Prometheus

GORM provides plugin `Prometheus` to collect `DBStats` and user-defined metrics

Check out [Prometheus](prometheus.html) for details

#### Naming Strategy

GORM allows users change the default naming conventions by overriding the default `NamingStrategy`, which is used to build `TableName`, `ColumnName`, `JoinTableName`, `RelationshipFKName`, `CheckerName`, `IndexName`, Check out [GORM Config](gorm_config.html) for details

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 NamingStrategy: schema.NamingStrategy{TablePrefix: "t_", SingularTable: true},
})
```

#### Logger

* Context support
* Customize/turn off the colors in the log
* Slow SQL log, default slow SQL time is 200ms
* Optimized the SQL log format so that it can be copied and executed in a database console

#### Transaction Mode

By default, all GORM write operations run inside a transaction to ensure data consistency, you can disable it during initialization to speed up write operations if it is not required

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
 SkipDefaultTransaction: true,
})
```

#### DataTypes (JSON as example)

GORM optimizes support for custom types, so you can define a struct to support all databases

The following takes JSON as an example (which supports SQLite, MySQL, Postgres, refer: https://github.com/go-gorm/datatypes/blob/master/json.go)

```go
import "gorm.io/datatypes"

type User struct {
 gorm.Model
 Name string
 Attributes datatypes.JSON
}

db.Create(&User{
 Name: "jinzhu",
 Attributes: datatypes.JSON([]byte(`{"name": "jinzhu", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`)),
}

// Query user having a role field in attributes
db.First(&user, datatypes.JSONQuery("attributes").HasKey("role"))
// Query user having orgs->orga field in attributes
db.First(&user, datatypes.JSONQuery("attributes").HasKey("orgs", "orga"))
```

#### Smart Select

GORM allows select specific fields with [`Select`](query.html), and in V2, GORM provides smart select mode if you are querying with a smaller struct

```go
type User struct {
 ID uint
 Name string
 Age int
 Gender string
 // hundreds of fields
}

type APIUser struct {
 ID uint
 Name string
}

// Select `id`, `name` automatically when query
db.Model(&User{}).Limit(10).Find(&APIUser{})
// SELECT `id`, `name` FROM `users` LIMIT 10
```

#### Associations Batch Mode

Association Mode supports batch data, e.g:

```go
// Find all roles for all users
db.Model(&users).Association("Role").Find(&roles)

// Delete User A from all user's team
db.Model(&users).Association("Team").Delete(&userA)

// Get unduplicated count of members in all user's team
db.Model(&users).Association("Team").Count()

// For `Append`, `Replace` with batch data, argument's length need to equal to data's length or will returns error
var users = []User{user1, user2, user3}
// e.g: we have 3 users, Append userA to user1's team, append userB to user2's team, append userA, userB and userC to user3's team
db.Model(&users).Association("Team").Append(&userA, &userB, &[]User{userA, userB, userC})
// Reset user1's team to userA，reset user2's team to userB, reset user3's team to userA, userB and userC
db.Model(&users).Association("Team").Replace(&userA, &userB, &[]User{userA, userB, userC})
```

#### Delete Associations when deleting

You are allowed to delete selected has one/has many/many2many relations with `Select` when deleting records, for example:

```go
// delete user's account when deleting user
db.Select("Account").Delete(&user)

// delete user's Orders, CreditCards relations when deleting user
db.Select("Orders", "CreditCards").Delete(&user)

// delete user's has one/many/many2many relations when deleting user
db.Select(clause.Associations).Delete(&user)

// delete user's account when deleting users
db.Select("Account").Delete(&users)
```

## Breaking Changes

We are trying to list big breaking changes or those changes can't be caught by the compilers, please create an issue or pull request [here](https://github.com/go-gorm/gorm.io) if you found any unlisted breaking changes

#### Tags

* GORM V2 prefer write tag name in `camelCase`, tags in `snake_case` won't works anymore, for example: `auto_increment`, `unique_index`, `polymorphic_value`, `embedded_prefix`, check out [Model Tags](models.html#tags)
* Tags used to specify foreign keys changed to `foreignKey`, `references`, check out [Associations Tags](associations.html#tags)
* Not support `sql` tag

#### Table Name

`TableName` will *not* allow dynamic table name anymore, the result of `TableName` will be cached for future

```go
func (User) TableName() string {
 return "t_user"
}
```

Please use `Scopes` for dynamic tables, for example:

```go
func UserTable(u *User) func(*gorm.DB) *gorm.DB {
 return func(db *gorm.DB) *gorm.DB {
 return db.Table("user_" + u.Role)
 }
}

db.Scopes(UserTable(&user)).Create(&user)
```

#### Creating and Deleting Tables requires the use of the Migrator

Previously tables could be created and dropped as follows:

```go
db.CreateTable(&MyTable{})
db.DropTable(&MyTable{})
```

Now you do the following:

```go
db.Migrator().CreateTable(&MyTable{})
db.Migrator().DropTable(&MyTable{})
```

#### Foreign Keys

A way of adding foreign key constraints was;

```go
db.Model(&MyTable{}).AddForeignKey("profile_id", "profiles(id)", "NO ACTION", "NO ACTION")
```

Now you add constraints as follows:

```go
db.Migrator().CreateConstraint(&Users{}, "Profiles")
db.Migrator().CreateConstraint(&Users{}, "fk_users_profiles")
```

which translates to the following sql code for postgres:

```sql
ALTER TABLE `Profiles` ADD CONSTRAINT `fk_users_profiles` FORIEGN KEY (`useres_id`) REFRENCES `users`(`id`))
```

#### Method Chain Safety/Goroutine Safety

To reduce GC allocs, GORM V2 will share `Statement` when using method chains, and will only create new `Statement` instances for new initialized `*gorm.DB` or after a `New Session Method`, to reuse a `*gorm.DB`, you need to make sure it just after a `New Session Method`, for example:

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

// Safe for new initialized *gorm.DB
for i := 0; i < 100; i++ {
 go db.Where(...).First(&user)
}

tx := db.Where("name = ?", "jinzhu")
// NOT Safe as reusing Statement
for i := 0; i < 100; i++ {
 go tx.Where(...).First(&user)
}

ctxDB := db.WithContext(ctx)
// Safe after a `New Session Method`
for i := 0; i < 100; i++ {
 go ctxDB.Where(...).First(&user)
}

ctxDB := db.Where("name = ?", "jinzhu").WithContext(ctx)
// Safe after a `New Session Method`
for i := 0; i < 100; i++ {
 go ctxDB.Where(...).First(&user) // `name = 'jinzhu'` will apply to the query
}

tx := db.Where("name = ?", "jinzhu").Session(&gorm.Session{})
// Safe after a `New Session Method`
for i := 0; i < 100; i++ {
 go tx.Where(...).First(&user) // `name = 'jinzhu'` will apply to the query
}
```

Check out [Method Chain](method_chaining.html) for details

#### Default Value

GORM V2 won't auto-reload default values created with database function after creating, checkout [Default Values](create.html#default_values) for details

#### Soft Delete

GORM V1 will enable soft delete if the model has a field named `DeletedAt`, in V2, you need to use `gorm.DeletedAt` for the model wants to enable the feature, e.g:

```go
type User struct {
 ID uint
 DeletedAt gorm.DeletedAt
}

type User struct {
 ID uint
 // field with different name
 Deleted gorm.DeletedAt
}
```

{% note warn %}
**NOTE:** `gorm.Model` is using `gorm.DeletedAt`, if you are embedding it, nothing needs to change
{% endnote %}

#### BlockGlobalUpdate

GORM V2 enabled `BlockGlobalUpdate` mode by default, to trigger a global update/delete, you have to use some conditions or use raw SQL or enable `AllowGlobalUpdate` mode, for example:

```go
db.Where("1 = 1").Delete(&User{})

db.Raw("delete from users")

db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})
```

#### ErrRecordNotFound

GORM V2 only returns `ErrRecordNotFound` when you are querying with methods `First`, `Last`, `Take` which is expected to return some result, and we have also removed method `RecordNotFound` in V2, please use `errors.Is` to check the error, e.g:

```go
err := db.First(&user).Error
errors.Is(err, gorm.ErrRecordNotFound)
```

#### Hooks Method

Before/After Create/Update/Save/Find/Delete must be defined as a method of type `func(tx *gorm.DB) error` in V2, which has unified interfaces like plugin callbacks, if defined as other types, a warning log will be printed and it won't take effect, check out [Hooks](hooks.html) for details

```go
func (user *User) BeforeCreate(tx *gorm.DB) error {
 // Modify current operation through tx.Statement, e.g:
 tx.Statement.Select("Name", "Age")
 tx.Statement.AddClause(clause.OnConflict{DoNothing: true})

 // Operations based on tx will runs inside same transaction without clauses of current one
 var role Role
 err := tx.First(&role, "name = ?", user.Role).Error
 // SELECT * FROM roles WHERE name = "admin"
 return err
}
```

#### Update Hooks support `Changed` to check fields changed or not

When updating with `Update`, `Updates`, You can use `Changed` method in Hooks `BeforeUpdate`, `BeforeSave` to check a field changed or not

```go
func (user *User) BeforeUpdate(tx *gorm.DB) error {
 if tx.Statement.Changed("Name", "Admin") { // if Name or Admin changed
 tx.Statement.SetColumn("Age", 18)
 }

 if tx.Statement.Changed() { // if any fields changed
 tx.Statement.SetColumn("Age", 18)
 }
 return nil
}

db.Model(&user).Update("Name", "Jinzhu") // update field `Name` to `Jinzhu`
db.Model(&user).Updates(map[string]interface{}{"name": "Jinzhu", "admin": false}) // update field `Name` to `Jinzhu`, `Admin` to false
db.Model(&user).Updates(User{Name: "Jinzhu", Admin: false}) // Update none zero fields when using struct as argument, will only update `Name` to `Jinzhu`

db.Model(&user).Select("Name", "Admin").Updates(User{Name: "Jinzhu"}) // update selected fields `Name`, `Admin`，`Admin` will be updated to zero value (false)
db.Model(&user).Select("Name", "Admin").Updates(map[string]interface{}{"Name": "Jinzhu"}) // update selected fields exists in the map, will only update field `Name` to `Jinzhu`

// Attention: `Changed` will only check the field value of `Update` / `Updates` equals `Model`'s field value, it returns true if not equal and the field will be saved
db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(map[string]interface{"name": "jinzhu2"}) // Changed("Name") => true
db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(map[string]interface{"name": "jinzhu"}) // Changed("Name") => false, `Name` not changed
db.Model(&User{ID: 1, Name: "jinzhu"}).Select("Admin").Updates(map[string]interface{"name": "jinzhu2", "admin": false}) // Changed("Name") => false, `Name` not selected to update

db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(User{Name: "jinzhu2"}) // Changed("Name") => true
db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(User{Name: "jinzhu"}) // Changed("Name") => false, `Name` not changed
db.Model(&User{ID: 1, Name: "jinzhu"}).Select("Admin").Updates(User{Name: "jinzhu2"}) // Changed("Name") => false, `Name` not selected to update
```

#### Plugins

Plugin callbacks also need be defined as a method of type `func(tx *gorm.DB) error`, check out [Write Plugins](write_plugins.html) for details

#### Updating with struct

When updating with struct, GORM V2 allows to use `Select` to select zero-value fields to update them, for example:

```go
db.Model(&user).Select("Role", "Age").Update(User{Name: "jinzhu", Role: "", Age: 0})
```

#### Associations

GORM V1 allows to use some settings to skip create/update associations, in V2, you can use `Select` to do the job, for example:

```go
db.Omit(clause.Associations).Create(&user)
db.Omit(clause.Associations).Save(&user)

db.Select("Company").Save(&user)
```

and GORM V2 doesn't allow preload with `Set("gorm:auto_preload", true)` anymore, you can use `Preload` with `clause.Associations`, e.g:

```go
// preload all associations
db.Preload(clause.Associations).Find(&users)
```

Also, checkout field permissions, which can be used to skip creating/updating associations globally

GORM V2 will use upsert to save associations when creating/updating a record, won't save full associations data anymore to protect your data from saving uncompleted data, for example:

```go
user := User{
 Name: "jinzhu",
 BillingAddress: Address{Address1: "Billing Address - Address 1"},
 ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
 Emails: []Email{
 {Email: "jinzhu@example.com"},
 {Email: "jinzhu-2@example.com"},
 },
 Languages: []Language{
 {Name: "ZH"},
 {Name: "EN"},
 },
}

db.Create(&user)
// BEGIN TRANSACTION;
// INSERT INTO "addresses" (address1) VALUES ("Billing Address - Address 1"), ("Shipping Address - Address 1") ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "users" (name,billing_address_id,shipping_address_id) VALUES ("jinzhu", 1, 2);
// INSERT INTO "emails" (user_id,email) VALUES (111, "jinzhu@example.com"), (111, "jinzhu-2@example.com") ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "languages" ("name") VALUES ('ZH'), ('EN') ON DUPLICATE KEY DO NOTHING;
// INSERT INTO "user_languages" ("user_id","language_id") VALUES (111, 1), (111, 2) ON DUPLICATE KEY DO NOTHING;
// COMMIT;
 ```

#### Join Table

In GORM V2, a `JoinTable` can be a full-featured model, with features like `Soft Delete`，`Hooks`, and define other fields, e.g:

```go
type Person struct {
 ID int
 Name string
 Addresses []Address `gorm:"many2many:person_addresses;"`
}

type Address struct {
 ID uint
 Name string
}

type PersonAddress struct {
 PersonID int
 AddressID int
 CreatedAt time.Time
 DeletedAt gorm.DeletedAt
}

func (PersonAddress) BeforeCreate(db *gorm.DB) error {
 // ...
}

// PersonAddress must defined all required foreign keys, or it will raise error
err := db.SetupJoinTable(&Person{}, "Addresses", &PersonAddress{})
```

After that, you could use normal GORM methods to operate the join table data, for example:

```go
var results []PersonAddress
db.Where("person_id = ?", person.ID).Find(&results)

db.Where("address_id = ?", address.ID).Delete(&PersonAddress{})

db.Create(&PersonAddress{PersonID: person.ID, AddressID: address.ID})
```

#### Count

Count only accepts `*int64` as the argument

#### Transactions

some transaction methods like `RollbackUnlessCommitted` removed, prefer to use method `Transaction` to wrap your transactions

```go
db.Transaction(func(tx *gorm.DB) error {
 // do some database operations in the transaction (use 'tx' from this point, not 'db')
 if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
 // return any error will rollback
 return err
 }

 if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
 return err
 }

 // return nil will commit the whole transaction
 return nil
})
```

Checkout [Transactions](transactions.html) for details

#### Migrator

* Migrator will create database foreign keys by default
* Migrator is more independent, many API renamed to provide better support for each database with unified API interfaces
* AutoMigrate will alter column's type if its size, precision, nullable changed
* Support Checker through tag `check`
* Enhanced tag setting for `index`

Checkout [Migration](migration.html) for details

```go
type UserIndex struct {
 Name string `gorm:"check:named_checker,(name <> 'jinzhu')"`
 Name2 string `gorm:"check:(age > 13)"`
 Name4 string `gorm:"index"`
 Name5 string `gorm:"index:idx_name,unique"`
 Name6 string `gorm:"index:,sort:desc,collate:utf8,type:btree,length:10,where:name3 != 'jinzhu'"`
}
```

## Happy Hacking!

<style>
li.toc-item { list-style: none; }
li.toc-item.toc-level-4 { display: none; }
</style>

---
title: Write Driver
layout: page
---

GORM offers built-in support for popular databases like `SQLite`, `MySQL`, `Postgres`, `SQLServer`, and `ClickHouse`. However, when you need to integrate GORM with databases that are not directly supported or have unique features, you can create a custom driver. This involves implementing the `Dialector` interface provided by GORM.

## Compatibility with MySQL or Postgres Dialects

For databases that closely resemble the behavior of `MySQL` or `Postgres`, you can often use the respective dialects directly. However, if your database significantly deviates from these dialects or offers additional features, developing a custom driver is recommended.

## Implementing the Dialector

The `Dialector` interface in GORM consists of methods that a database driver must implement to facilitate communication between the database and GORM. Let's break down the key methods:

```go
type Dialector interface {
 Name() string // Returns the name of the database dialect
 Initialize(*DB) error // Initializes the database connection
 Migrator(db *DB) Migrator // Provides the database migration tool
 DataTypeOf(*schema.Field) string // Determines the data type for a schema field
 DefaultValueOf(*schema.Field) clause.Expression // Provides the default value for a schema field
 BindVarTo(writer clause.Writer, stmt *Statement, v interface{}) // Handles variable binding in SQL statements
 QuoteTo(clause.Writer, string) // Manages quoting of identifiers
 Explain(sql string, vars ...interface{}) string // Formats SQL statements with variables
}
```

Each method in this interface serves a crucial role in how GORM interacts with the database, from establishing connections to handling queries and migrations.

### Nested Transaction Support

If your database supports savepoints, you can implement the `SavePointerDialectorInterface` to get the `Nested Transaction Support` and `SavePoint` support.

```go
type SavePointerDialectorInterface interface {
	SavePoint(tx *DB, name string) error // Saves a savepoint within a transaction
	RollbackTo(tx *DB, name string) error // Rolls back a transaction to the specified savepoint
}
```

By implementing these methods, you enable support for savepoints and nested transactions, offering advanced transaction management capabilities.

### Custom Clause Builders

Defining custom clause builders in GORM allows you to extend the query capabilities for specific database operations. In this example, we'll go through the steps to define a custom clause builder for the "LIMIT" clause, which may have database-specific behavior.

- **Step 1: Define a Custom Clause Builder Function**:

To create a custom clause builder, you need to define a function that adheres to the `clause.ClauseBuilder` interface. This function will be responsible for constructing the SQL clause for a specific operation. In our example, we'll create a custom "LIMIT" clause builder.

Here's the basic structure of a custom "LIMIT" clause builder function:

```go
func MyCustomLimitBuilder(c clause.Clause, builder clause.Builder) {
 if limit, ok := c.Expression.(clause.Limit); ok {
 // Handle the "LIMIT" clause logic here
 // You can access the limit values using limit.Limit and limit.Offset
 builder.WriteString("MYLIMIT")
 }
}
```

- The function takes two parameters: `c` of type `clause.Clause` and `builder` of type `clause.Builder`.
- Inside the function, we check if the `c.Expression` is a `clause.Limit`. If it is, we proceed to handle the "LIMIT" clause logic.

Replace `MYLIMIT` with the actual SQL logic for your database. This is where you can implement database-specific behavior for the "LIMIT" clause.

- **Step 2: Register the Custom Clause Builder**:

To make your custom "LIMIT" clause builder available to GORM, register it with the `db.ClauseBuilders` map, typically during driver initialization. Here's how to register the custom "LIMIT" clause builder:

```go
func (d *MyDialector) Initialize(db *gorm.DB) error {
 // Register the custom "LIMIT" clause builder
 db.ClauseBuilders["LIMIT"] = MyCustomLimitBuilder

 //...
}
```

In this code, we use the key `"LIMIT"` to register our custom clause builder in the `db.ClauseBuilders` map, associating our custom builder with the "LIMIT" clause.

- **Step 3: Use the Custom Clause Builder**:

After registering the custom clause builder, GORM will call it when generating SQL statements that involve the "LIMIT" clause. You can use your custom logic to generate the SQL clause as needed.

Here's an example of how you can use the custom "LIMIT" clause builder in a GORM query:

```go
query := db.Model(&User{})

// Apply the custom "LIMIT" clause using the Limit method
query = query.Limit(10) // You can also provide an offset, e.g., query.Limit(10).Offset(5)

// Execute the query
result := query.Find(&results)
// SQL: SELECT * FROM users MYLIMIT
```

In this example, we use the Limit method with GORM, and behind the scenes, our custom "LIMIT" clause builder (MyCustomLimitBuilder) will be invoked to handle the generation of the "LIMIT" clause.

For inspiration and guidance, examining the [MySQL Driver](https://github.com/go-gorm/mysql) can be helpful. This driver demonstrates how the `Dialector` interface is implemented to suit the specific needs of the MySQL database.

---
title: Write Plugins
layout: page
---

## Callbacks

GORM leverages `Callbacks` to power its core functionalities. These callbacks provide hooks for various database operations like `Create`, `Query`, `Update`, `Delete`, `Row`, and `Raw`, allowing for extensive customization of GORM's behavior.

Callbacks are registered at the global `*gorm.DB` level, not on a session basis. This means if you need different callback behaviors, you should initialize a separate `*gorm.DB` instance.

### Registering a Callback

You can register a callback for specific operations. For example, to add a custom image cropping functionality:

```go
func cropImage(db *gorm.DB) {
 if db.Statement.Schema != nil {
 // crop image fields and upload them to CDN, dummy code
 for _, field := range db.Statement.Schema.Fields {
 switch db.Statement.ReflectValue.Kind() {
 case reflect.Slice, reflect.Array:
 for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
 // Get value from field
 if fieldValue, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue.Index(i)); !isZero {
 if crop, ok := fieldValue.(CropInterface); ok {
 crop.Crop()
 }
 }
 }
 case reflect.Struct:
 // Get value from field
 if fieldValue, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue); !isZero {
 if crop, ok := fieldValue.(CropInterface); ok {
 crop.Crop()
 }
 }

 // Set value to field
 err := field.Set(db.Statement.Context, db.Statement.ReflectValue, "newValue")
 }
 }

 // All fields for current model
 db.Statement.Schema.Fields

 // All primary key fields for current model
 db.Statement.Schema.PrimaryFields

 // Prioritized primary key field: field with DB name `id` or the first defined primary key
 db.Statement.Schema.PrioritizedPrimaryField

 // All relationships for current model
 db.Statement.Schema.Relationships

 // Find field with field name or db name
 field := db.Statement.Schema.LookUpField("Name")

 // processing
 }
}

// Register the callback for the Create operation
db.Callback().Create().Register("crop_image", cropImage)
```

### Deleting a Callback

If a callback is no longer needed, it can be removed:

```go
// Remove the 'gorm:create' callback from Create operations
db.Callback().Create().Remove("gorm:create")
```

### Replacing a Callback

Callbacks with the same name can be replaced with a new function:

```go
// Replace the 'gorm:create' callback with a new function
db.Callback().Create().Replace("gorm:create", newCreateFunction)
```

### Ordering Callbacks

Callbacks can be registered with specific orders to ensure they execute at the right time in the operation lifecycle.

```go
// Register to execute before the 'gorm:create' callback
db.Callback().Create().Before("gorm:create").Register("update_created_at", updateCreated)

// Register to execute after the 'gorm:create' callback
db.Callback().Create().After("gorm:create").Register("update_created_at", updateCreated)

// Register to execute after the 'gorm:query' callback
db.Callback().Query().After("gorm:query").Register("my_plugin:after_query", afterQuery)

// Register to execute after the 'gorm:delete' callback
db.Callback().Delete().After("gorm:delete").Register("my_plugin:after_delete", afterDelete)

// Register to execute before the 'gorm:update' callback
db.Callback().Update().Before("gorm:update").Register("my_plugin:before_update", beforeUpdate)

// Register to execute before 'gorm:create' and after 'gorm:before_create'
db.Callback().Create().Before("gorm:create").After("gorm:before_create").Register("my_plugin:before_create", beforeCreate)

// Register to execute before any other callbacks
db.Callback().Create().Before("*").Register("update_created_at", updateCreated)

// Register to execute after any other callbacks
db.Callback().Create().After("*").Register("update_created_at", updateCreated)
```

### Predefined Callbacks

GORM comes with a set of predefined callbacks that drive its standard features. It's recommended to review these [defined callbacks](https://github.com/go-gorm/gorm/blob/master/callbacks/callbacks.go) before creating custom plugins or additional callback functions.

## Plugins

GORM's plugin system allows for easy extensibility and customization of its core functionalities, enhancing your application's capabilities while maintaining a modular architecture.

### The `Plugin` Interface

To create a plugin for GORM, you need to define a struct that implements the `Plugin` interface:

```go
type Plugin interface {
 Name() string
 Initialize(*gorm.DB) error
}
```

- **`Name` Method**: Returns a unique string identifier for the plugin.
- **`Initialize` Method**: Contains the logic to set up the plugin. This method is called when the plugin is registered with GORM for the first time.

### Registering a Plugin

Once your plugin conforms to the `Plugin` interface, you can register it with a GORM instance:

```go
// Example of registering a plugin
db.Use(MyCustomPlugin{})
```

### Accessing Registered Plugins

After a plugin is registered, it is stored in GORM's configuration. You can access registered plugins via the `Plugins` map:

```go
// Access a registered plugin by its name
plugin := db.Config.Plugins[pluginName]
```

### Practical Example

An example of a GORM plugin is the Prometheus plugin, which integrates Prometheus monitoring with GORM:

```go
// Registering the Prometheus plugin
db.Use(prometheus.New(prometheus.Config{
 // Configuration options here
}))
```

[Prometheus plugin documentation](prometheus.html) provides detailed information on its implementation and usage.
