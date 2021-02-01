import React, { useState } from 'react'

import { postJSON } from '../../utils/api'
import { formatAsDate, formatAsMoney } from '../../utils/format'
import LoaderComponent from '../common/LoaderComponent'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'
import CardActivitiesTable from '../tables/CardActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const LinkerFlowPage = ({ match }) => {
    const splitwiseExpenseUUID = match.params.splitwiseExpenseUUID

    const [splitwiseExpense, setSplitwiseExpense] = useState(null)
    const [linkLoading, setLinkLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    const linkCardActivityToExpense = (cardActivityUUID) => {
        setLinkLoading(true)
        postJSON(`/api/card_activities/${cardActivityUUID}/link/${splitwiseExpenseUUID}`)
            .then(data => {
                if (data.message === 'success') {
                    window.history.back()
                }
            })
            .catch(e => setErrorMessage(e.message))
            .finally(() => setLinkLoading(false))
    }

    const linkAccountActivityToExpense = (accountActivityUUID) => {
        setLinkLoading(true)
        postJSON(`/api/account_activities/${accountActivityUUID}/link/${splitwiseExpenseUUID}`)
            .then(data => {
                if (data.message === 'success') {
                    window.history.back()
                }
            })
            .catch(e => setErrorMessage(e.message))
            .finally(() => setLinkLoading(false))
    }

    let linksDiv = null
    if (splitwiseExpense !== null) {
        linksDiv = <div>
            <h2>Possible Links</h2>

            <h3>Card Activity Links</h3>
            <CardActivitiesTable
                data={splitwiseExpense.card_activity_links}
                hideFilters={true}
                extraColumns={['link']}
                extraCustomRenders={{
                    'transaction_date': (data) => <div style={{
                        color: formatAsDate(data['transaction_date']) === formatAsDate(splitwiseExpense['date']) ? 'green' : null
                    }}>{formatAsDate(data['transaction_date'])}</div>,
                    'post_date': (data) => <div style={{
                        color: formatAsDate(data['post_date']) === formatAsDate(splitwiseExpense['date']) ? 'green' : null
                    }}>{formatAsDate(data['post_date'])}</div>,
                    'amount': (data) => <div style={{
                        color: Math.abs(data['amount']) === Math.abs(splitwiseExpense['amount_paid']) ? 'green' : null
                    }}>{formatAsMoney(data['amount'])}</div>,
                    'link': (data) => <input
                        type='button'
                        onClick={() => linkCardActivityToExpense(data['uuid'])}
                        value='Link'
                    />,
                }}
            />

            <h3>Account Activity Links</h3>
            <AccountActivitiesTable
                data={splitwiseExpense.account_activity_links}
                hideFilters={true}
                extraColumns={['link']}
                extraCustomRenders={{
                    'posting_date': (data) => <div style={{
                        color: formatAsDate(data['posting_date']) === formatAsDate(splitwiseExpense['date']) ? 'green' : null
                    }}>{formatAsDate(data['posting_date'])}</div>,
                    'amount': (data) => <div style={{
                        color: Math.abs(data['amount']) === Math.abs(splitwiseExpense['amount_paid']) ? 'green' : null
                    }}>{formatAsMoney(data['amount'])}</div>,
                    'link': (data) => <input
                        type='button'
                        onClick={() => linkAccountActivityToExpense(data['uuid'])}
                        value='Link'
                    />
                }}
            />
        </div>
    }

    return (
        <div>
            <h1>Splitwise Expense Linking Flow</h1>
            <LoaderComponent
                path={`/api/splitwise_expenses/${splitwiseExpenseUUID}/for_linking`}
                parentLoading={linkLoading}
                setData={setSplitwiseExpense}
            />
            <SplitwiseExpenseTable data={splitwiseExpense !== null ? [splitwiseExpense] : []} hideFilters={true} />

            {linksDiv}
        </div>
    )
}

export default LinkerFlowPage
