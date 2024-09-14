// Create Category
​
POST: /category
​
input: {
    name: string
    profile_ids?: string[]
}
output: {
    category_id: string
}

///////
// Update category
​
PATCH: /category/${id}
​
input: {
    name: string
    profile_ids?: string[]
}
​
output: {
    name: string
    profile_ids?: string[]
}
////////
// get all category
​
GET: /category?page=1&size=10&keyword=
// page: default: 1
// size: default: 10
​
output: {
    data: {
        id: string
        name: string
        number_profile: number
        created_at: string | Date
        updated_at: string | Date
    }[],
    current: number
    total: number
}

///////////
/// del category
DELETE: /category/${id}
​
output: {
    status: boolean
}

/////
/// remove profile id in category
PATCH: /category/${id}/${profile_id}
​
output: {
    status: boolean
}
//////
/// get detail category
​
enum EPrivacy {
    PUBLIC = "public",
    PRIVATE = "private",
    PROTECTED = 'protected'
}
​
type TProfileValue = {
    headers: string[] // [expected_return, standard_deviation]
    values: string[]
}
​
type TProfile = {
    id: string
    name: string
    charts: TProfileValue
    privacy: EPrivacy
    author_id: string
    total_return: number
    last_updated: string | Date
}
​
GET: /category/${category_id}?page=1&size=10
​
output: {
    id: string
    name: string
    profiles: TProfile[]
    current: number
    total: number
}
​
