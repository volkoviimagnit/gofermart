package response

type OrderStatus string

const AccrualStatusRegistered OrderStatus = "REGISTERED"
const AccrualStatusInvalid OrderStatus = "INVALID"
const AccrualStatusProcessing OrderStatus = "PROCESSING"
const AccrualStatusProcessed OrderStatus = "PROCESSED"
