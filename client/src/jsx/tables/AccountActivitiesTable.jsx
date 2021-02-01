import React from 'react'
import { Link } from 'react-router-dom'

import { formatAsDate, formatAsMoney } from '../../utils/format'
import { dateComparator } from '../../utils/sorting'
import ExtendableTable from './ExtendableTable'

const AccountActivitiesTable = (props) => {
    return <ExtendableTable
        columns={[
            'uuid',
            'details',
            'posting_date',
            'description',
            'amount',
            'type',
        ]}
        customRenders={{
            'uuid': (data) =>
                <Link to={'/accounts/' + data['account_uuid'] + '/activities/' + data['uuid']}>{data['uuid']}</Link>,
            'posting_date': (data) => formatAsDate(data['posting_date']),
            'amount': (data) => formatAsMoney(data['amount']),
        }}
        initialSortColumn='posting_date'
        customSortComparators={{
            'posting_date': dateComparator
        }}
        {...props}
    />
}

export default AccountActivitiesTable
