export interface User {
  id: string
  username: string
  email: string
  password: string
  created: string
  logged_in: string
}

export interface Item {
  id: string
  name: string
  description: string
  image_url: string
  url: string
  currency: string
}

export interface ItemWithPrice {
  id: string
  name: string
  description: string
  image_url: string
  url: string
  currency: string
  time: string
  price: number
  available: boolean
}

export interface UserItem {
  user_id: string
  item_id: string
}

export interface ItemPrice {
  item_id: string
  time: string
  price: number
  available: boolean
}

export interface Subscription {
  user_id: string
  item_id: string
  email: string
  target_price: string
}
