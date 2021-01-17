import Moment from 'moment'

export function dateComparator(a, b) {
    return Moment(a).valueOf() - Moment(b).valueOf()
}

export function defaultSort(a, b) {
    if (a < b) {
        return -1
    } else if (a > b) {
        return 1
    } else {
        return 0
    }
}
