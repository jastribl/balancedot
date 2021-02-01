import React, { useEffect, useState } from 'react'

import { get, postJSON } from '../../utils/api'
import { formatAsDate, formatAsMoney } from '../../utils/format'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'
import Table from '../common/Table'
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
        setExpenseLoading(true)
        get(`/api/splitwise_expenses/${splitwiseExpenseUUID}`)
            .then(splitwiseExpenseResponse => {
                setSplitwiseExpense(splitwiseExpenseResponse)
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setExpenseLoading(false)
            })

        setCardLinksLoading(true)
        get(`/api/card_activities/for_link/${splitwiseExpenseUUID}`)
            .then(cardLinksResponse => {
                setCardLinks(cardLinksResponse)
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setCardLinksLoading(false)
            })

        setAccountLinksLoading(true)
        get(`/api/account_activities/for_link/${splitwiseExpenseUUID}`)
            .then(accountLinksResponse => {
                setAccountLinks(accountLinksResponse)
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setAccountLinksLoading(false)
            })
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
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setLinkLoading(false)
            })
    }

    const linkAccountActivityToExpense = (accountActivityUUID) => {
        setLinkLoading(true)
        postJSON(`/api/account_activities/${accountActivityUUID}/link/${splitwiseExpenseUUID}`)
            .then(data => {
                if (data.message === 'success') {
                    window.history.back()
                }
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setLinkLoading(false)
            })
    }

    return (
        <div>
            <Spinner visible={expenseLoading || cardLinksLoading || accountLinksLoading || linkLoading} />
            <h1>Splitwise Expense Linking Flow</h1>
            <ErrorRow message={errorMessage} />
            <SplitwiseExpenseTable data={splitwiseExpense ? [splitwiseExpense] : []} hideFilters={true} />
            <h2>Possible Links</h2>

            <h3>Card Links</h3>
            <div>
                <Table
                    rowKey='uuid'
                    rows={cardLinks}
                    columns={['transaction_date', 'post_date', 'description', 'category', 'type', 'amount', 'link']}
                    customRenders={{
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
                        />
                    }}
                    hideFilters={true}
                />

                <h3>Account Links</h3>
                <Table
                    rowKey='uuid'
                    rows={accountLinks}
                    columns={['uuid', 'details', 'posting_date', 'description', 'amount', 'type', 'link']}
                    customRenders={{
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
                    hideFilters={true}
                />
            </div>
        </div>
    )
}

export default LinkerFlowPage
