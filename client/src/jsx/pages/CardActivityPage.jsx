import React, { useState } from 'react'

import LoaderComponent from '../common/LoaderComponent'
import CardActivitiesTable from '../tables/CardActivitiesTable'
import SplitwiseExpenseTable from '../tables/SplitwiseExpenseTable'

const CardActivityPage = ({ match }) => {
    const cardActivityUUID = match.params.cardActivityUUID

    const [cardActivity, setCardActivity] = useState(null)

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
            <h1>Card Activity {cardActivityUUID} ({cardActivity?.description}) </h1>
            <h2>For card {card ? (card.last_four + " (" + card.description + ")") : null}</h2>
            <LoaderComponent
                path={`/api/card_activities/${cardActivityUUID}`}
                parentLoading={false}
                setData={setCardActivity}
            />
            <CardActivitiesTable
                data={cardActivity ? [cardActivity] : []}
                hideFilters={true}
            />
            {splitwiseExpenseTable}
        </div>
    )
}

export default CardActivityPage
