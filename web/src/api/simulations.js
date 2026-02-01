import { api } from 'boot/axios'

/**
 * Fetch all available simulations
 * @returns {Promise} Axios response promise
 */
export function fetchSimulations() {
  return api.get('/v1/public/simulations')
}

/**
 * Get a single simulation by ID
 * @param {string} id - Simulation ID
 * @returns {Promise} Axios response promise
 */
export function getSimulation(id) {
  return api.get(`/v1/public/simulations/${id}`)
}

/**
 * Start a simulation run
 * @param {string} simulationId - ID of the simulation to run
 * @param {Object} parameters - Parameter values for the run
 * @returns {Promise} Axios response promise
 */
export function runSimulation(simulationId, parameters) {
  return api.post('/v1/public/simulations/run', {
    simulation_id: simulationId,
    parameters,
  })
}

/**
 * Get the current run status
 * @returns {Promise} Axios response promise
 */
export function getRunStatus() {
  return api.get('/v1/public/simulations/run')
}

/**
 * Fetch the output of the current run as binary data
 * @returns {Promise} Axios response promise
 */
export function fetchRunOutput(config = {}) {
  return api.get('/v1/public/simulations/output?format=json', {
    responseType: 'arraybuffer',
    ...config,
  })
}
