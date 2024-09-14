type TProfileValue = {
    headers: string[] // [expected_return, standard_deviation]
    values: string[]
}
​
type Author = {
    id: string
    name: string
}
​
type TProfile = {
    id: string
    name: string
    charts: TProfileValue
    privacy: EPrivacy
    author: Author
    total_return: number
    updated_at: number
    created_at: number
}
​
​
type TPerformanceKey = '1w' | '1M' | '3M' | '6M' | '1Y' | '3Y' | 'All' // Edit??
​
type TPerformanceValue = {
     value: number
     status: boolean // true: tăng, false: giảm. Giá tăng giảm phụ thuộc vào giá close.
}
​
type TCharts = {
    name: string
    portfolio_profile: number
    sp_asx: number
}
​

////////////
// Get Profile byUserId
GET: /porfolio-profile?page=1&size=10&keyword=
​
output: {
    data: TProfile[]
    current: number
    total: number
}
///////////
// Get Detail Profile
GET: /porfolio-profile/${profile_id}
​
output: {
    id: string
    name: string
    privacy: EPrivacy
    author: Author
​
    category?: {
        id: string
        name: string
    }
    create_updated: string | Date
    updated_by: string | Date
    advisor?: {
        id: string
        code: string
    }[] | null
    branch?: {
        id: string
        code: string
    }[] | null
    organization?: {
        id: string
        code: string
    }[] | null
    privacy: EPrivacy
    number_linked_accounts?: number // all profile is linked // TODO:
    performance: Record<TPerformanceKey, TPerformanceValue>
    statistics: {
        headers: ['Portfolio', 'Portfolio Profile name', 'S&P/ASX 200'],
        values: [string, number, number][]
    },
    allocation: {
        headers: string[]
        values: number[]
    }[],
    profile_detail: {
        headers: string[]
        values: (string | number)[]
    }
    histories: {
        headers: string[]
        values: (string | number)[]
    }
}
​
//////////////
// GEt  charts Opt Profile
​
GET: /porfolio-profile/${profile_id}/opt?start={'SOY' | '1D' | '1M' | '3M' | '1Y' | '5Y'}
​
​
output: {
    data: TCharts[]
}
///////////
// Copy Profile
POST: /porfolio-profile/${profile_id}/copy
​
output: {
    status: boolean
}
///////////
// api get all Linked profile
POST: /porfolio-profile/{profile_id}/linked_accounts?keyword=
​
output: {
    data: {
        id: string
        name: string
        charts: TCharts[]
    }[],
    current: number
    total: number
}
//////////////
// api link account to profile
​
POST: /porfolio-profile/{profile_id}/linked_accounts
​
input: {
    account_ids: string[]
}
output: {
    status: boolean
}
​
///////////////
// get Branch code, Organisation Code, Advisor Code
​
GET: /porfolio-profile/{branch|organization|advisor}?keyword= // /porfolio-profile/branch
// default limit 10
​
output: {
    data: {
        id: string
        code: string
    }[]
}
///////////////
// create profile
POST: /porfolio-profile
​
input: {
    category_id?: string
    name: string
    organization_id?: string[]
    branch_id?: string[]
    advisor_id?: string[]
    assets?: {
        ticker_id: string
        allocation: number
        price?: string
    }[],
    privacy: EPrivacy
}
​
output: {
    profile_id: string
}
//////////////
// update Profile
​
PATCH: /porfolio-profile/${profile_id}
input: {
    category_id?: string
    name?: string
    organization_id?: string[]
    branch_id?: string[]
    advisor_id?: string[]
    assets?: {
        ticker_id: string
        allocation: number
        price?: string
    }[],
    privacy: EPrivacy
}

output: {
  status: boolean
}
​
//////////////
// del profile
DELETE: /porfolio-profile/${profile_id}
​
​
output: {
    status: boolean
}
