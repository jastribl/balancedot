import Moment from 'moment'
import React from 'react'

import { formatAsDate, formatAsMoney } from '../../../utils/format'
import AccountActivitiesTable from '../../tables/AccountActivitiesTable'
import CardActivitiesTable from '../../tables/CardActivitiesTable'

const SplitwiseLinkingSection = ({ splitwiseExpense, handleLinking }) => {
    if (splitwiseExpense === null) {
        return null
    }

    const getDateDiffDays = (a, b) => Math.abs(Moment.utc(a).diff(Moment.utc(b), 'days'))
    const getMoneyDiffCents = (a, b) =>
        Math.round(Math.abs(Math.abs(Math.abs(a) - Math.abs(b))).toFixed(2) * 100)

    const getDateStyle = (itemDate, expenseDate) => {
        switch (getDateDiffDays(itemDate, expenseDate)) {
            case 0:
                return {
                    color: '#00ff00'
                }
            default:
                return {}
        }
    }

    const getMoneyStyle = (itemMoney, expenseMoney) => {
        if ((itemMoney > 0) !== (expenseMoney > 0)) {
            const diffCents = getMoneyDiffCents(itemMoney, expenseMoney)
            if (diffCents === 0) {
                return {
                    color: '#00ff00'
                }
            } else if (diffCents <= 3) {
                return {
                    color: '#0000ff'
                }
            }
        }
        return null
    }

    const getCardActivityDistanceNumber = (data) => {
        const transactionDateDiff = getDateDiffDays(data['transaction_date'], splitwiseExpense['date'])
        const postDateDiff = getDateDiffDays(data['post_date'], splitwiseExpense['date'])
        const averageDateDiff = (transactionDateDiff + postDateDiff) / 2
        const moneyDiffCents = getMoneyDiffCents(data['amount'], splitwiseExpense['amount_paid'])

        return Math.min(moneyDiffCents, 5) + Math.min(averageDateDiff, 3)
        // return averageDateDiff + Math.floor((moneyDiffCents * 1.0) / 5)
    }

    const getAccountActivityDistanceNumber = (data) => {
        const postingDateDiff = getDateDiffDays(data['posting_date'], splitwiseExpense['date'])
        const moneyDiffCents = getMoneyDiffCents(data['amount'], splitwiseExpense['amount_paid'])

        return Math.max(moneyDiffCents, 5) + Math.max(postingDateDiff, 3)
        // return postingDateDiff + Math.floor((moneyDiff * 1.0) / 5)
    }

    if (splitwiseExpense !== null && (
        splitwiseExpense.card_activity_links || splitwiseExpense.account_activity_links
    )) {
        let cardLinksDiv = null
        if (splitwiseExpense.card_activity_links?.length > 0) {
            cardLinksDiv = <div>
                <h3>Card Activity Links</h3>
                <CardActivitiesTable
                    data={splitwiseExpense.card_activity_links}
                    extraColumns={['diff', 'link']}
                    extraCustomRenders={{
                        'transaction_date': (data) => <div style={
                            getDateStyle(data['transaction_date'], splitwiseExpense['date'])
                        }>{formatAsDate(data['transaction_date'])}</div>,
                        'post_date': (data) => <div style={
                            getDateStyle(data['post_date'], splitwiseExpense['date'])
                        }>{formatAsDate(data['post_date'])}</div>,
                        'amount': (data) => <div style={
                            getMoneyStyle(data['amount'], splitwiseExpense['amount_paid'])
                        }>{formatAsMoney(data['amount'])}</div>,
                        'diff': (data) => getCardActivityDistanceNumber(data),
                        'link': (data) => <input
                            type='button'
                            onClick={() => handleLinking('card_activities', 'link', data['uuid'])}
                            value='Link'
                        />,
                    }}
                    initialSortColumn={null}
                    customSortFunctionOverride={(a, b) => {
                        return getCardActivityDistanceNumber(a) - getCardActivityDistanceNumber(b)
                    }}
                />
            </div>
        }
        let accountLinksDiv = null
        if (splitwiseExpense.account_activity_links?.length > 0) {
            accountLinksDiv = <div>
                <h3>Account Activity Links</h3>
                <AccountActivitiesTable
                    data={splitwiseExpense.account_activity_links}
                    extraColumns={['diff', 'link']}
                    extraCustomRenders={{
                        'posting_date': (data) => <div style={
                            getDateStyle(data['posting_date'], splitwiseExpense['date'])
                        }>{formatAsDate(data['posting_date'])}</div>,
                        'amount': (data) => <div style={
                            getMoneyStyle(data['amount'], splitwiseExpense['amount_paid'])
                        }>{formatAsMoney(data['amount'])}</div>,
                        'diff': (data) => getAccountActivityDistanceNumber(data),
                        'link': (data) => <input
                            type='button'
                            onClick={() => handleLinking('account_activities', 'link', data['uuid'])}
                            value='Link'
                        />
                    }}
                    initialSortColumn={null}
                    customSortFunctionOverride={(a, b) => {
                        return getAccountActivityDistanceNumber(a) - getAccountActivityDistanceNumber(b)
                    }}
                />
            </div>
        }
        if (cardLinksDiv !== null || accountLinksDiv !== null) {
            return <div>
                <h2>Possible Links</h2>
                {cardLinksDiv}
                {accountLinksDiv}
            </div>
        } else {
            return <h2>No Links Found :(</h2>
        }
    }

    return null
}

export default SplitwiseLinkingSection
