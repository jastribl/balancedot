import React, { useEffect, useState } from 'react'

import { getWithHandling, postJSON } from '../../utils/api'
import { formatAsDate, formatAsMoney } from '../../utils/format'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'
import CardActivitiesTable from '../tables/CardActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const LinkerFlowPage = ({ match }) => {
    const splitwiseExpenseUUID = match.params.splitwiseExpenseUUID

    const [splitwiseExpense, setSplitwiseExpense] = useState(null)
    const [cardLinks, setCardLinks] = useState(null)
    const [accountLinks, setAccountLinks] = useState(null)
    const [expenseLoading, setExpenseLoading] = useState(false)
    const [cardLinksLoading, setCardLinksLoading] = useState(false)
    const [accountLinksLoading, setAccountLinksLoading] = useState(false)
    const [linkLoading, setLinkLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    useEffect(() => {
        getWithHandling(
            `/api/splitwise_expenses/${splitwiseExpenseUUID}`,
            setSplitwiseExpense,
            setErrorMessage,
            setExpenseLoading
        )

        getWithHandling(
            `/api/card_activities/for_link/${splitwiseExpenseUUID}`,
            setCardLinks,
            setErrorMessage,
            setCardLinksLoading
        )

        getWithHandling(
            `/api/account_activities/for_link/${splitwiseExpenseUUID}`,
            setAccountLinks,
            setErrorMessage,
            setAccountLinksLoading
        )
    }, [
        setExpenseLoading,
        setCardLinksLoading,
        setAccountLinksLoading,
        setSplitwiseExpense,
        setCardLinks,
        setAccountLinks,
        setErrorMessage,
    ])

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
                data={cardLinks}
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
                data={accountLinks}
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
            <Spinner visible={expenseLoading || cardLinksLoading || accountLinksLoading || linkLoading} />
            <h1>Splitwise Expense Linking Flow</h1>
            <ErrorRow message={errorMessage} />
            <SplitwiseExpenseTable data={splitwiseExpense !== null ? [splitwiseExpense] : []} hideFilters={true} />

            {linksDiv}
        </div>
    )
}

export default LinkerFlowPage
