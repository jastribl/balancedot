import Moment from 'moment'

export function formatAsMoney(amount, currencyCode = 'USD') {
    let currencySymbol = '?'
    switch (currencyCode) {
        case 'USD':
            currencySymbol = '$'
            break
        default:
            currencySymbol = `${currencyCode}?`
    }
    return (amount < 0 ? '-' : '') + currencySymbol + Math.abs(amount).toFixed(2)
}

export function formatAsDate(date) {
    if (date !== null) {
        return Moment.utc(date).format('YYYY-MM-DD')
    }
    return '????-??-??'
}