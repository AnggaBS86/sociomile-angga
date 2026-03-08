const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {}),
      ...(options.headers || {})
    },
    method: options.method || 'GET',
    body: options.body ? JSON.stringify(options.body) : undefined
  })

  const text = await response.text()
  const payload = text ? JSON.parse(text) : null

  if (!response.ok) {
    const error = new Error(payload?.message || 'Request failed')
    error.status = response.status
    error.payload = payload
    throw error
  }

  return payload
}

export function login(email, password) {
  return request('/auth/login', {
    method: 'POST',
    body: { email, password }
  })
}

export function listConversations(token) {
  return request('/conversations', { token })
}

export function getConversationDetail(token, conversationId) {
  return request(`/conversations/${conversationId}`, { token })
}

export function replyConversation(token, conversationId, message) {
  return request(`/conversations/${conversationId}/messages`, {
    method: 'POST',
    token,
    body: { message }
  })
}

export function escalateConversation(token, conversationId, body) {
  return request(`/conversations/${conversationId}/escalate`, {
    method: 'POST',
    token,
    body
  })
}

export function listTickets(token) {
  return request('/tickets', { token })
}

export function updateTicketStatus(token, ticketId, status) {
  return request(`/tickets/${ticketId}/status`, {
    method: 'PATCH',
    token,
    body: { status }
  })
}
