import React, { useEffect, useState } from 'react'

import { get } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'
import SplitwiseLoginCheck from '../SplitwiseLoginCheck'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'
import CardActivitiesTable from '../tables/CardActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const SplitwiseExpensePage = ({ match }) => {
    const splitwiseExpenseUUID = match.params.splitwiseExpenseUUID

    const [splitwiseExpense, setSplitwiseExpense] = useState(null)
    const [pageLoading, setPageLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    const refreshSplitwiseExpense = () => {
        setPageLoading(true)
        get(`/api/splitwise_expenses/${splitwiseExpenseUUID}`)
            .then(splitwiseExpenseResponse => {
                setSplitwiseExpense(splitwiseExpenseResponse)
            })
            .catch(e => {
                setErrorMessage(e.message)
            })
            .finally(() => {
                setPageLoading(false)
            })
    }


    useEffect(() => {
        refreshSplitwiseExpense()
    }, [setPageLoading, setSplitwiseExpense, setErrorMessage])

    let cardActivitiesTable = null
    if (splitwiseExpense?.card_activities !== null && splitwiseExpense?.card_activities.length > 0) {
        cardActivitiesTable = <div>
            <h3>Card Activities</h3>
            <CardActivitiesTable data={splitwiseExpense?.card_activities} hideFilters={true} />
        </div>
    }


    let accountActivitiesTable = null
    if (splitwiseExpense?.accountivities !== null && splitwiseExpense?.account_activities.length > 0) {
        accountActivitiesTable = <div>
            <h3>Account Activities</h3>
            <AccountActivitiesTable data={splitwiseExpense?.account_activities} hideFilters={true} />
        </div>
    }

    return (
        <div>
            <Spinner visible={pageLoading} />
            <h1>Splitwise Expense {splitwiseExpenseUUID} ({splitwiseExpense?.description})</h1>
            <ErrorRow message={errorMessage} />
            <SplitwiseLoginCheck>
                <SplitwiseExpenseTable data={splitwiseExpense ? [splitwiseExpense] : []} />
                {cardActivitiesTable}
                {accountActivitiesTable}
            </SplitwiseLoginCheck>
        </div>
    )
}

export default SplitwiseExpensePage
