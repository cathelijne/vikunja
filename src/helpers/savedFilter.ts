import type {IList} from '@/modelTypes/IList'

export function getSavedFilterIdFromListId(listId: IList['id']) {
	let filterId = listId * -1 - 1
	// FilterIds from listIds are always positive
	if (filterId < 0) {
		filterId = 0
	}
	return filterId
}

export function isSavedFilter(list: IList) {
	return getSavedFilterIdFromListId(list.id) > 0
}