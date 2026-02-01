import { api } from 'boot/axios'

/**
 * Initialize a new client session
 * @returns {Promise} Axios response promise
 */
export function initializeClient() {
  return api.post('/v1/public/client/init')
}

/**
 * Fetch an existing client by ID (uses header from axios config)
 * @returns {Promise} Axios response promise
 */
export function fetchClient() {
  return api.get('/v1/public/client')
}
