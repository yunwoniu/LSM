package sortTree

import "sync"

type node struct {
	kv    *Kv
	left  *node
	right *node
}

type SortTree struct {
	lock sync.RWMutex
	root *node
	size int
}

func (s *SortTree) Set(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	newNode := &node{kv: &Kv{Key: key, Val: val}}
	if s.root == nil {
		s.size++
		s.root = newNode
		return
	}
	cur := s.root
	for cur != nil {
		if cur.kv.Key < key {
			if cur.left == nil {
				cur.left = newNode
				s.size++
				return
			}
			cur = cur.left
		} else if cur.kv.Key == key {
			cur.kv = &Kv{Key: key, Val: val}
			return
		} else {
			if cur.right == nil {
				cur.right = newNode
				s.size++
				return
			}
			cur = cur.right
		}
	}
	panic(any("this is bug"))
	return
}

func (s *SortTree) Get(key string) (string, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	cur := s.root
	for cur != nil {
		if cur.kv.Key < key {
			if cur.left == nil {
				return "", false
			}
			cur = cur.left
		} else if cur.kv.Key == key {
			if cur.kv.Delete {
				return "", false
			}
			return cur.kv.Val, true
		} else {
			if cur.right == nil {
				return "", false
			}
			cur = cur.right
		}
	}
	return "", false
}

func (s *SortTree) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	cur := s.root
	for cur != nil {
		if cur.kv.Key < key {
			if cur.left == nil { //没找到就直接返回
				return
			}
			cur = cur.left
		} else if cur.kv.Key == key { //找到了就把它置为删除标记置为true
			cur.kv.Delete = true
			return
		} else {
			if cur.right == nil { //没找到就直接返回
				return
			}
			cur = cur.right
		}
	}
}

func (s *SortTree) GetCount() int {
	return s.size
}

func (s *SortTree) GetSortKv() []*Kv {
	st := make([]*node, 0, s.GetCount())
	if s.root == nil {
		return nil
	}
	var ret []*Kv
	st = append(st, s.root)
	for len(st) != 0 {
		top := st[len(st)-1]
		st = st[:len(st)-1]
		if top != nil { //中序遍历
			if top.right != nil {
				st = append(st, top.right)
			}
			st = append(st, top)
			st = append(st, nil)
			if top.left != nil {
				st = append(st, top.left)
			}
		} else {
			node := st[len(st)-1]
			ret = append(ret, node.kv)
			st = st[:len(st)-1]
		}
	}
	return ret
}
