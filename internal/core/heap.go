package core

type BidsHeap []*OrderItem
type AsksHeap []*OrderItem

func (h BidsHeap) Len() int {
	return len(h)
}

func (h BidsHeap) Less(i, j int) bool {
	return h[i].Price < h[j].Price
}

func (h BidsHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *BidsHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	item.index = -1
	*h = old[0 : n-1]
	return item
}

func (h *BidsHeap) Push(x any) {
	item := x.(*OrderItem)
	item.index = len(*h)
	*h = append(*h, item)
}

func (h BidsHeap) Peek() *OrderItem {
	if len(h) == 0 {
		return nil
	}

	return h[0]
}

func (h AsksHeap) Len() int {
	return len(h)
}

func (h AsksHeap) Less(i, j int) bool {
	return h[i].Price > h[j].Price
}

func (h AsksHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *AsksHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	item.index = -1
	*h = old[0 : n-1]
	return item
}

func (h *AsksHeap) Push(x any) {
	item := x.(*OrderItem)
	*h = append(*h, item)
	item.index = len(*h)
}

func (h AsksHeap) Peek() *OrderItem {
	if len(h) == 0 {
		return nil
	}

	return h[0]
}
