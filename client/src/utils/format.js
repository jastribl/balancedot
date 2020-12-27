export function formatAsMoney(amount, currencyCode = 'USD') {
    let currencySymbol = '?'
    switch (currencyCode) {
        case 'USD':
            currencySymbol = '$';
            break;
        default:
            currencySymbol = `${currencyCode}?`
    }
    return (amount < 0 ? '-' : '') + currencySymbol + Math.abs(amount).toFixed(2)
}
