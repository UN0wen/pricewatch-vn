import { AxiosInstance } from './axios'
import { ItemWithPrice, UserItem } from './models'

type CreateUserItemPayload = {
  item: {
    id: string
  }
}

// Returns an array of items
export async function getAllUserItems() {
  try {
    const response = await AxiosInstance.get('/user/items')
    const data = response.data
    if (Array.isArray(data)) {
      const items = data.map(d => <ItemWithPrice>d.item_with_price)

      return items
    }
    return
  } catch (err) {
    console.log(err.response.data.error)
  }
}

export async function createUserItem(payload: CreateUserItemPayload) {
  const requestOptions = {
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  }
  try {
    const response = await AxiosInstance.post(`/user/item`, payload, requestOptions)
    const data = response.data
    if (data) {
      return <UserItem>data.user_item
    }
    return
  } catch (err) {
    console.log(err.response.data.error)
  }
}

export async function checkUserItem(id: string) : Promise<boolean> {

  try {
    await AxiosInstance.get(`/user/item/${id}`)
    
    // 200 code
    return true
  } catch (err) {
    console.log(err.response.data.error)
    return false
  }
}
