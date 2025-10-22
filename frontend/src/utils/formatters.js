import { format, parseISO } from 'date-fns'

export const formatCurrency = (amount, currency = 'VND') => {
  if (currency === 'VND') {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount)
  }
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount)
}

export const formatDate = (date) => {
  const dateObj = typeof date === 'string' ? parseISO(date) : date
  return format(dateObj, 'MMM dd, yyyy')
}

export const formatDateTime = (date) => {
  const dateObj = typeof date === 'string' ? parseISO(date) : date
  return format(dateObj, 'MMM dd, yyyy HH:mm')
}

export const formatTime = (date) => {
  const dateObj = typeof date === 'string' ? parseISO(date) : date
  return format(dateObj, 'HH:mm')
}

export const formatNumber = (num, decimals = 0) => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }).format(num)
}

export const formatPercentage = (value, decimals = 1) => {
  return new Intl.NumberFormat('en-US', {
    style: 'percent',
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }).format(value / 100)
}

export const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

export const formatDuration = (seconds) => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60
  if (hours > 0) return `${hours}h ${minutes}m ${secs}s`
  if (minutes > 0) return `${minutes}m ${secs}s`
  return `${secs}s`
}

export const truncateText = (text, maxLength = 50) => {
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

export const capitalize = (text) => {
  return text.charAt(0).toUpperCase() + text.slice(1).toLowerCase()
}

export const formatTransactionType = (type) => {
  const types = { income: 'Income', expense: 'Expense', transfer: 'Transfer' }
  return types[type] || capitalize(type)
}

export const formatGoalType = (type) => {
  const types = {
    savings: 'Savings',
    debt_payment: 'Debt Payment',
    investment: 'Investment',
    purchase: 'Purchase',
    other: 'Other',
  }
  return types[type] || capitalize(type)
}

export const formatPriority = (priority) => {
  const priorities = { low: 'Low', medium: 'Medium', high: 'High', urgent: 'Urgent' }
  return priorities[priority] || capitalize(priority)
}

export const formatPeriod = (period) => {
  const periods = { weekly: 'Weekly', monthly: 'Monthly', yearly: 'Yearly' }
  return periods[period] || capitalize(period)
}

export function formatVietnamDate(isoString, withTime = false) {
  if (!isoString) return "";

  const date = new Date(isoString);

  // Chuyển sang giờ Việt Nam thủ công
  const utc = date.getTime() + date.getTimezoneOffset() * 60000;
  const vnDate = new Date(utc + 7 * 60 * 60000);

  const day = String(vnDate.getDate()).padStart(2, "0");
  const month = String(vnDate.getMonth() + 1).padStart(2, "0");
  const year = vnDate.getFullYear();

  let formatted = `${day}/${month}/${year}`;

  if (withTime) {
    const hours = String(vnDate.getHours()).padStart(2, "0");
    const minutes = String(vnDate.getMinutes()).padStart(2, "0");
    const seconds = String(vnDate.getSeconds()).padStart(2, "0");
    formatted += ` ${hours}:${minutes}:${seconds}`;
  }

  return formatted;
}

