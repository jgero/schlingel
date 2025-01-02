// @generated automatically by Diesel CLI.

diesel::table! {
    resource_groups (id) {
        id -> Uuid,
        name -> Text,
        resource_group_id -> Nullable<Uuid>,
    }
}

diesel::table! {
    resources (id) {
        id -> Uuid,
        name -> Text,
        hash -> Nullable<Bytea>,
        resource_group_id -> Nullable<Uuid>,
    }
}

diesel::joinable!(resources -> resource_groups (resource_group_id));

diesel::allow_tables_to_appear_in_same_query!(resource_groups, resources,);
