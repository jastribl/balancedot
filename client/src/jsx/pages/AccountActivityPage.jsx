import React, { useEffect, useState } from 'react'

import { get } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const AccountActivityPage = ({ match }) => {
    const accountActivityUUID = match.params.accountActivityUUID

    const [accountActivity, setAccountActivity] = useState(null)
    const [accountActivityLoading, setAccountActivityLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)

    const refreshAccountActivity = () => {
        setAccountActivityLoading(true)
        get(`/api/account_activities/${accountActivityUUID}`)
            .then(accountActivityResponse => setAccountActivity(accountActivityResponse))
            .catch(e => setErrorMessage(e.message))
            .finally(() => setAccountActivityLoading(false))
    }

    useEffect(() => {
        refreshAccountActivity()
    }, [setAccountActivity])

    let splitwiseExpenseTable = null
    if (accountActivity?.splitwise_expenses !== null && accountActivity?.splitwise_expenses.length > 0) {
        splitwiseExpenseTable = <div>
            <h3>Splitwise Expenses</h3>
            <SplitwiseExpenseTable data={accountActivity?.splitwise_expenses} hideFilters={true} />
        </div>
    }

    const account = accountActivity?.account

    return (
        <div>
            <Spinner visible={accountActivityLoading} />
            <h1>Account Activity {accountActivityUUID} ({accountActivity?.description}) </h1>
            <h2>For account {account ? (account.last_four + " (" + account.description + ")") : null}</h2>
            <ErrorRow message={errorMessage} />
            <AccountActivitiesTable
                data={accountActivity ? [accountActivity] : []}
                hideFilters={true}
            />
            {splitwiseExpenseTable}
        </div>
    )
}

export default AccountActivityPage
