use diesel::prelude::*;

#[derive(Queryable, Selectable, Identifiable, Associations, Debug, PartialEq)]
#[diesel(belongs_to(ResourceGroup))]
#[diesel(table_name = crate::db::resource_groups)]
pub struct ResourceGroup {
    pub id: uuid::Uuid,
    pub name: String,
    pub resource_group_id: Option<uuid::Uuid>,
}

#[derive(Queryable, Selectable, Identifiable, Associations, Debug, PartialEq)]
#[diesel(belongs_to(ResourceGroup))]
#[diesel(table_name = crate::db::resources)]
pub struct Resource {
    pub id: uuid::Uuid,
    pub name: String,
    pub resource_group_id: Option<uuid::Uuid>,
}
