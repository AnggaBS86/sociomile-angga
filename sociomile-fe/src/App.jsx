import { useEffect, useMemo, useState } from 'react'
import {
  escalateConversation,
  getConversationDetail,
  listConversations,
  listTickets,
  login,
  replyConversation,
  updateTicketStatus
} from './api'

const ticketStatuses = ['open', 'assigned', 'closed']

function extractMessage(error) {
  if (error?.payload?.errors?.length) {
    return error.payload.errors.map((item) => `${item.field}: ${item.message}`).join(', ')
  }
  return error?.message || 'Unexpected error'
}

function StatusBadge({ value }) {
  return <span className={`badge badge-${value || 'default'}`}>{value || 'unknown'}</span>
}

function EmptyState({ title, body }) {
  return (
    <div className="empty-state">
      <h3>{title}</h3>
      <p>{body}</p>
    </div>
  )
}

export default function App() {
  const [token, setToken] = useState(localStorage.getItem('sociomile_token') || '')
  const [user, setUser] = useState(() => {
    const raw = localStorage.getItem('sociomile_user')
    return raw ? JSON.parse(raw) : null
  })
  const [email, setEmail] = useState('angga@email.com')
  const [password, setPassword] = useState('123456')
  const [authError, setAuthError] = useState('')
  const [loadingAuth, setLoadingAuth] = useState(false)

  const [conversations, setConversations] = useState([])
  const [tickets, setTickets] = useState([])
  const [selectedConversationId, setSelectedConversationId] = useState(null)
  const [selectedConversation, setSelectedConversation] = useState(null)
  const [conversationLoading, setConversationLoading] = useState(false)
  const [replyText, setReplyText] = useState('')
  const [replyLoading, setReplyLoading] = useState(false)
  const [escalateForm, setEscalateForm] = useState({
    title: 'Escalated issue',
    description: 'Need internal follow up',
    priority: 'high'
  })
  const [escalateLoading, setEscalateLoading] = useState(false)
  const [statusMessage, setStatusMessage] = useState('')
  const [errorMessage, setErrorMessage] = useState('')
  const [ticketUpdateLoadingId, setTicketUpdateLoadingId] = useState(null)

  const isAuthenticated = Boolean(token)
  const isAdmin = user?.role === 'admin'
  const isAgent = user?.role === 'agent'

  const selectedConversationSummary = useMemo(
    () => conversations.find((item) => item.id === selectedConversationId) || null,
    [conversations, selectedConversationId]
  )

  async function hydrateDashboard(authToken, nextSelectionId) {
    setErrorMessage('')
    setStatusMessage('')
    try {
      const [conversationResponse, ticketResponse] = await Promise.all([
        listConversations(authToken),
        listTickets(authToken)
      ])
      const nextConversations = conversationResponse?.data || []
      const nextTickets = ticketResponse?.data || []
      setConversations(nextConversations)
      setTickets(nextTickets)

      const preferredId = nextSelectionId || nextConversations[0]?.id || null
      setSelectedConversationId(preferredId)
      if (preferredId) {
        await loadConversationDetail(authToken, preferredId)
      } else {
        setSelectedConversation(null)
      }
    } catch (error) {
      setErrorMessage(extractMessage(error))
    }
  }

  async function loadConversationDetail(authToken, conversationId) {
    setConversationLoading(true)
    try {
      const response = await getConversationDetail(authToken, conversationId)
      setSelectedConversation(response?.data || null)
    } catch (error) {
      setErrorMessage(extractMessage(error))
    } finally {
      setConversationLoading(false)
    }
  }

  useEffect(() => {
    if (!token) {
      return
    }
    hydrateDashboard(token)
  }, [token])

  async function handleLogin(event) {
    event.preventDefault()
    setLoadingAuth(true)
    setAuthError('')
    try {
      const response = await login(email, password)
      const nextToken = response.data.token
      const nextUser = response.data.user
      localStorage.setItem('sociomile_token', nextToken)
      localStorage.setItem('sociomile_user', JSON.stringify(nextUser))
      setToken(nextToken)
      setUser(nextUser)
    } catch (error) {
      setAuthError(extractMessage(error))
    } finally {
      setLoadingAuth(false)
    }
  }

  function handleLogout() {
    localStorage.removeItem('sociomile_token')
    localStorage.removeItem('sociomile_user')
    setToken('')
    setUser(null)
    setConversations([])
    setTickets([])
    setSelectedConversationId(null)
    setSelectedConversation(null)
    setStatusMessage('')
    setErrorMessage('')
  }

  async function handleSelectConversation(conversationId) {
    setSelectedConversationId(conversationId)
    await loadConversationDetail(token, conversationId)
  }

  async function handleReply(event) {
    event.preventDefault()
    if (!selectedConversationId || !replyText.trim()) {
      return
    }
    setReplyLoading(true)
    setErrorMessage('')
    setStatusMessage('')
    try {
      await replyConversation(token, selectedConversationId, replyText.trim())
      setReplyText('')
      setStatusMessage('Reply sent')
      await hydrateDashboard(token, selectedConversationId)
    } catch (error) {
      setErrorMessage(extractMessage(error))
    } finally {
      setReplyLoading(false)
    }
  }

  async function handleEscalate(event) {
    event.preventDefault()
    if (!selectedConversationId) {
      return
    }
    setEscalateLoading(true)
    setErrorMessage('')
    setStatusMessage('')
    try {
      const response = await escalateConversation(token, selectedConversationId, escalateForm)
      setStatusMessage(response?.message || 'Escalation queued')
      await hydrateDashboard(token, selectedConversationId)
    } catch (error) {
      setErrorMessage(extractMessage(error))
    } finally {
      setEscalateLoading(false)
    }
  }

  async function handleTicketStatus(ticketId, status) {
    setTicketUpdateLoadingId(ticketId)
    setErrorMessage('')
    setStatusMessage('')
    try {
      const response = await updateTicketStatus(token, ticketId, status)
      setStatusMessage(response?.message || 'Ticket updated')
      await hydrateDashboard(token, selectedConversationId)
    } catch (error) {
      setErrorMessage(extractMessage(error))
    } finally {
      setTicketUpdateLoadingId(null)
    }
  }

  if (!isAuthenticated) {
    return (
      <main className="auth-layout">
        <section className="auth-panel brand-panel">
          <p className="eyebrow">Sociomile 2.0</p>
          <h1>Inbox, escalation, and ticket operations in one workspace.</h1>
          <p className="description">
            This frontend is wired to the Echo backend. Use the seeded account to log in and test the workflow.
          </p>
        </section>
        <section className="auth-panel form-panel">
          <form onSubmit={handleLogin} className="auth-form">
            <h2>Sign in</h2>
            <label>
              <span>Email</span>
              <input value={email} onChange={(event) => setEmail(event.target.value)} type="email" />
            </label>
            <label>
              <span>Password</span>
              <input value={password} onChange={(event) => setPassword(event.target.value)} type="password" />
            </label>
            {authError ? <p className="error-text">{authError}</p> : null}
            <button type="submit" disabled={loadingAuth}>
              {loadingAuth ? 'Signing in...' : 'Login'}
            </button>
          </form>
        </section>
      </main>
    )
  }

  return (
    <main className="dashboard-shell">
      <header className="topbar">
        <div>
          <p className="eyebrow">Sociomile Inbox</p>
          <h1>Operations Console</h1>
        </div>
        <div className="topbar-actions">
          <div className="user-chip">
            <strong>{user?.email}</strong>
            <span>{user?.role} | tenant {user?.tenant_id}</span>
          </div>
          <button className="ghost-button" onClick={() => hydrateDashboard(token, selectedConversationId)}>
            Refresh
          </button>
          <button className="ghost-button" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </header>

      {statusMessage ? <p className="status-banner success-banner">{statusMessage}</p> : null}
      {errorMessage ? <p className="status-banner error-banner">{errorMessage}</p> : null}

      <section className="dashboard-grid">
        <aside className="panel list-panel">
          <div className="panel-heading">
            <h2>Conversations</h2>
            <span>{conversations.length}</span>
          </div>
          {conversations.length === 0 ? (
            <EmptyState title="No conversations" body="Push a webhook event from Postman or the backend to create one." />
          ) : (
            <div className="list-stack">
              {conversations.map((conversation) => (
                <button
                  key={conversation.id}
                  className={`conversation-item ${selectedConversationId === conversation.id ? 'active' : ''}`}
                  onClick={() => handleSelectConversation(conversation.id)}
                >
                  <div>
                    <strong>Conversation #{conversation.id}</strong>
                    <p>Customer #{conversation.customer_id}</p>
                  </div>
                  <StatusBadge value={conversation.status} />
                </button>
              ))}
            </div>
          )}
        </aside>

        <section className="panel detail-panel">
          <div className="panel-heading">
            <h2>Conversation Detail</h2>
            {selectedConversationSummary ? <StatusBadge value={selectedConversationSummary.status} /> : null}
          </div>
          {conversationLoading ? <p>Loading conversation...</p> : null}
          {!conversationLoading && !selectedConversation ? (
            <EmptyState title="Select a conversation" body="Choose a conversation from the left panel to inspect messages and reply." />
          ) : null}
          {!conversationLoading && selectedConversation ? (
            <>
              <div className="detail-meta">
                <span>ID: {selectedConversation.conversation.id}</span>
                <span>Customer: {selectedConversation.conversation.customer_id}</span>
                <span>
                  Assigned: {selectedConversation.conversation.assigned_agent_id || 'unassigned'}
                </span>
              </div>
              <div className="message-timeline">
                {selectedConversation.messages.length === 0 ? (
                  <EmptyState title="No messages" body="This conversation has no messages yet." />
                ) : (
                  selectedConversation.messages.map((message) => (
                    <article key={message.id} className={`message-card sender-${message.sender_type}`}>
                      <div className="message-meta">
                        <strong>{message.sender_type}</strong>
                        <span>{message.created_at || 'now'}</span>
                      </div>
                      <p>{message.message}</p>
                    </article>
                  ))
                )}
              </div>

              <form className="composer" onSubmit={handleReply}>
                <label>
                  <span>Reply</span>
                  <textarea
                    rows="4"
                    value={replyText}
                    onChange={(event) => setReplyText(event.target.value)}
                    placeholder="Write a response to the customer"
                  />
                </label>
                <button type="submit" disabled={replyLoading}>
                  {replyLoading ? 'Sending...' : 'Send Reply'}
                </button>
              </form>

              {isAgent ? (
                <form className="escalation-form" onSubmit={handleEscalate}>
                  <div className="panel-heading compact">
                    <h3>Escalate to Ticket</h3>
                    <span>Agent only</span>
                  </div>
                  <label>
                    <span>Title</span>
                    <input
                      value={escalateForm.title}
                      onChange={(event) => setEscalateForm((prev) => ({ ...prev, title: event.target.value }))}
                    />
                  </label>
                  <label>
                    <span>Description</span>
                    <textarea
                      rows="3"
                      value={escalateForm.description}
                      onChange={(event) => setEscalateForm((prev) => ({ ...prev, description: event.target.value }))}
                    />
                  </label>
                  <label>
                    <span>Priority</span>
                    <select
                      value={escalateForm.priority}
                      onChange={(event) => setEscalateForm((prev) => ({ ...prev, priority: event.target.value }))}
                    >
                      <option value="low">low</option>
                      <option value="medium">medium</option>
                      <option value="high">high</option>
                    </select>
                  </label>
                  <button type="submit" disabled={escalateLoading}>
                    {escalateLoading ? 'Queueing...' : 'Escalate'}
                  </button>
                </form>
              ) : null}
            </>
          ) : null}
        </section>

        <aside className="panel ticket-panel">
          <div className="panel-heading">
            <h2>Tickets</h2>
            <span>{tickets.length}</span>
          </div>
          {tickets.length === 0 ? (
            <EmptyState title="No tickets" body="Escalated conversations will appear here." />
          ) : (
            <div className="ticket-list">
              {tickets.map((ticket) => (
                <article className="ticket-card" key={ticket.id}>
                  <div className="ticket-header">
                    <div>
                      <strong>{ticket.title}</strong>
                      <p>Conversation #{ticket.conversation_id}</p>
                    </div>
                    <StatusBadge value={ticket.status} />
                  </div>
                  <p className="ticket-description">{ticket.description || 'No description'}</p>
                  <div className="ticket-footer">
                    <span>Priority: {ticket.priority}</span>
                    {isAdmin ? (
                      <select
                        value={ticket.status}
                        onChange={(event) => handleTicketStatus(ticket.id, event.target.value)}
                        disabled={ticketUpdateLoadingId === ticket.id}
                      >
                        {ticketStatuses.map((status) => (
                          <option value={status} key={status}>
                            {status}
                          </option>
                        ))}
                      </select>
                    ) : null}
                  </div>
                </article>
              ))}
            </div>
          )}
        </aside>
      </section>
    </main>
  )
}
