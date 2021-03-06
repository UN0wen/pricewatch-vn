import { AxiosInstance } from './axios'
import { Item, ItemPrice, ItemWithPrice } from './models'

type CreateItemPayload = {
  item: {
    url: string
  }
}

// Returns an array of items with prices
export async function getAllItems() {
  try {
    const response = await AxiosInstance.get('/items/prices')
    const data = response.data
    if (Array.isArray(data)) {
      const items = data.map((d) => <ItemWithPrice>d.item_with_price)
      return items
    }
    return null
  } catch (err) {
    console.log(err.response.data.error)
    return null
  }
}

// Returns an array of items with prices
export async function search(params: URLSearchParams) {
  try {
    const response = await AxiosInstance.get('/items/search', { params })
    const data = response.data
    if (Array.isArray(data)) {
      const items = data.map((d) => <ItemWithPrice>d.item_with_price)
      return items
    }
    return null
  } catch (err) {
    console.log(err.response.data.error)
    return null
  }
}

// Returns an item with its price
export async function getItem(id: string) {
  try {
    const response = await AxiosInstance.get(`/item/${id}`)
    const data = response.data
    if (data.item_with_price) {
      return <ItemWithPrice>data.item_with_price
    }
    return null
  } catch (err) {
    console.log(err.response.data.error)
    return null
  }
}

// Returns an item's latest price
export async function getItemPrice(id: string) {
  try {
    const response = await AxiosInstance.get(`/item/${id}/price`)
    const data = response.data
    if (data.item) {
      return <ItemPrice>data.price
    }
    return null
  } catch (err) {
    console.log(err.response.data.error)
    return null
  }
}

// Returns all item's prices
export async function getItemPrices(id: string) {
  try {
    const response = await AxiosInstance.get(`/item/${id}/prices`)
    const data = response.data
    if (Array.isArray(data)) {
      const itemPrices = data.map((d) => <ItemPrice>d.price)
      return itemPrices
    }
    return null
  } catch (err) {
    console.log(err.response.data.error)
    return null
  }
}

export async function createItem(item: any) {
  const payload = {
    item: item
  }
  const requestOptions = {
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  }
  try {
    const response = await AxiosInstance.post(`/item`, payload, requestOptions)
    const data = response.data
    if (data) {
      return <Item>data.item
    }
    return null
  } catch (err) {
    console.log(err.response.data.error)
    return null
  }
}

export async function checkURL(payload: CreateItemPayload): Promise<boolean> {
  const requestOptions = {
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  }
  try {
    await AxiosInstance.post(`/item/validate`, payload, requestOptions)

    return true
  } catch (err) {
    console.log(err.response.data.error)
    return false
  }
}

export async function getItemFromURL(payload: CreateItemPayload): Promise<any> {
  const requestOptions = {
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  }
  try {
    const response = await AxiosInstance.post(`/item/url`, payload, requestOptions)

    const data = response.data
    if (data) {
      return <Item>data.item
    }

    return 400
  } catch (err) {
    return Number(err.response.status)
  }
}