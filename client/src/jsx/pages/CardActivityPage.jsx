import React, { useEffect, useState } from 'react'

import { getWithHandling } from '../../utils/api'
import ErrorRow from '../common/ErrorRow'
import Spinner from '../common/Spinner'
import CardActivitiesTable from '../tables/CardActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const CardActivityPage = ({ match }) => {
    const cardActivityUUID = match.params.cardActivityUUID

    const [cardActivity, setCardActivity] = useState(null)
    const [cardActivityLoading, setCardActivityLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState(null)


    useEffect(() => {
        getWithHandling(
            `/api/card_activities/${cardActivityUUID}`,
            setCardActivity,
            setErrorMessage,
            setCardActivityLoading
        )
    }, [
        setCardActivity,
        setErrorMessage,
        setCardActivityLoading,
    ])

    let splitwiseExpenseTable = null
    if (cardActivity?.splitwise_expenses !== null && cardActivity?.splitwise_expenses.length > 0) {
        splitwiseExpenseTable = <div>
            <h3>Splitwise Expenses</h3>
            <SplitwiseExpenseTable data={cardActivity?.splitwise_expenses} hideFilters={true} />
        </div>
    }

    const card = cardActivity?.card

    return (
        <div>
            <Spinner visible={cardActivityLoading} />
            <h1>Card Activity {cardActivityUUID} ({cardActivity?.description}) </h1>
            <h2>For card {card ? (card.last_four + " (" + card.description + ")") : null}</h2>
            <ErrorRow message={errorMessage} />
            <CardActivitiesTable
                data={cardActivity ? [cardActivity] : []}
                hideFilters={true}
            />
            {splitwiseExpenseTable}
        </div>
    )
}

export default CardActivityPage
