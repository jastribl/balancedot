import React, { useEffect, useState } from 'react'

import { get } from '../../utils/api'
import AccountActivitiesTable from '../tables/AccountActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const AccountActivityPage = ({ match }) => {
    const accountUUID = match.params.accountUUID
    const accountActivityUUID = match.params.accountActivityUUID

    const [account, setAccount] = useState(null)
    const [accountActivity, setAccountActivity] = useState(null)

    const refreshAccount = () => {
        get(`/api/accounts/${accountUUID}`)
            .then(accountResponse => setAccount(accountResponse))
    }

    const refreshAccountActivity = () => {
        get(`/api/account_activities/${accountActivityUUID}`)
            .then(accountActivityResponse => setAccountActivity(accountActivityResponse))
    }

    useEffect(() => {
        refreshAccount()
        refreshAccountActivity()
    }, [setAccount, setAccountActivity])

    let splitwiseExpenseTable = null
    if (accountActivity?.splitwise_expenses !== null && accountActivity?.splitwise_expenses.length > 0) {
        splitwiseExpenseTable = <div>
            <h3>Splitwise Expenses</h3>
            <SplitwiseExpenseTable data={accountActivity?.splitwise_expenses} hideFilters={true} />
        </div>
    }

    return (
        <div>
            <h1>Account Activity {accountActivityUUID} ({accountActivity?.description}) </h1>
            <h2>For account {account ? (account.last_four + " (" + account.description + ")") : null}</h2>
            <AccountActivitiesTable
                data={accountActivity ? [accountActivity] : []}
                hideFilters={true}
            />
            {splitwiseExpenseTable}
        </div>
    )
}

export default AccountActivityPage
