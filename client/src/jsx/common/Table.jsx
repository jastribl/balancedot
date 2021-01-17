import React, { useState } from 'react'

import { snakeToSentenceCase } from '../../utils/strings'

const Table = ({ rowKey, columns, rows, customRenders, initialSortColumn, initialSortInverse, customSortComparators }) => {
    customRenders ??= {}
    customSortComparators ??= {}

    if (!rows) {
        return <div />
    }

    const [sortColumn, setSortColumn] = useState(initialSortColumn)
    const [sortInverse, setSortInverse] = useState(initialSortInverse ?? false)

    const onHeaderClick = (header_name) => {
        if (sortInverse) {
            setSortColumn(null)
            setSortInverse(false)
        } else if (sortColumn === header_name) {
            setSortInverse(!sortInverse)
        } else {
            setSortColumn(header_name)
            setSortInverse(false)
        }
    }

    let toRender = rows.slice()
    if (sortColumn) {
        toRender.sort(((a, b) => customSortComparators[sortColumn](a[sortColumn], b[sortColumn])) ??
            ((a, b) => {
                if (a[sortColumn] < b[sortColumn]) {
                    return -1
                } else if (a[sortColumn] > b[sortColumn]) {
                    return 1
                } else {
                    return 0
                }
            }))
    }
    if (sortInverse) {
        toRender.reverse()
    }

    return (
        <div>
            <table className='styled-table'>
                <thead>
                    <tr>
                        {columns.map(key =>
                            <th
                                key={key}
                                onClick={() => onHeaderClick(key)}
                            >{snakeToSentenceCase(key)}{(key === sortColumn ? (sortInverse ? " ↑" : " ↓") : "")}</th>
                        )}
                    </tr>
                </thead>
                <tbody>
                    {toRender.map((row, i) =>
                        <tr key={row[rowKey]}>{
                            columns.map(key =>
                                <td key={key}>{
                                    key in customRenders ? customRenders[key](row) : row[key]
                                }</td>
                            )
                        }</tr>
                    )}
                </tbody>
            </table>
        </div>
    )
}

export default Table