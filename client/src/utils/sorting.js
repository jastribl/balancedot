import Moment from 'moment'

export function dateComparator(a, b) {
    return Moment(a).valueOf() - Moment(b).valueOf()
}
