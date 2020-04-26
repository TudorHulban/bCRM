/*
package models Contains application models.

How security works:
A. For user information fetching and management
1. there is app admin that can create, fetch, delete and update users. secu level: 4
2. there is group admin that has rights for group as app admin for app. secu level: 3
3. there is team admin. team admin may be normal user for another team. secu level: 2
3. there is team managers. manager is just a label. same security rights as normal user. no secu level for now.
4. there is normal user. normal user may be in different teams. secu level: 0

(secu level 1 is reserved, maybe something comes out).

Modeling:
1. define teams table
2. define teams security table: ID, Team ID, User ID, type of access.
*/
package models
