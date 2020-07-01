package page

import (
	"fmt"
	"strconv"
	"sync"
)

// Page represents a single page.
type Page struct {
	sync.RWMutex
	Name          string             `json:"name"`
	Controls      map[string]Control `json:"controls"`
	nextControlID int
	clients       map[string]*Client
	clientsMutex  sync.RWMutex
}

// NewPage creates a new instance of Page.
func NewPage(name string) (*Page, error) {
	p := &Page{}
	p.Name = name
	p.Controls = make(map[string]Control)
	p.AddControl(NewControl("Page", "", p.NextControlID()))
	p.clients = make(map[string]*Client)
	return p, Pages().Add(p)
}

// NextControlID returns the next auto-generated control ID
func (page *Page) NextControlID() string {
	page.Lock()
	defer page.Unlock()
	nextID := strconv.Itoa(page.nextControlID)
	page.nextControlID++
	return nextID
}

// AddControl adds a control to a page
func (page *Page) AddControl(ctl Control) error {
	// find parent
	parentID := ctl.ParentID()
	if parentID != "" {
		page.RLock()
		parentCtl, ok := page.Controls[parentID]
		page.RUnlock()

		if !ok {
			return fmt.Errorf("parent control with id '%s' not found", parentID)
		}

		// update parent's childIds
		page.Lock()
		parentCtl.AddChildID(ctl.ID())
		page.Unlock()
	}

	page.Lock()
	defer page.Unlock()
	page.Controls[ctl.ID()] = ctl
	return nil
}

func (page *Page) registerClient(client *Client) {
	page.clientsMutex.Lock()
	defer page.clientsMutex.Unlock()
	page.clients[client.id] = client
}

func (page *Page) unregisterClient(client *Client) {
	page.clientsMutex.Lock()
	defer page.clientsMutex.Unlock()
	delete(page.clients, client.id)
}

// func (p Page) MarshalJSON() ([]byte, error) {
// 	var tmp struct {
// 		Name     string             `json:"name"`
// 		Controls map[string]Control `json:"controls1"`
// 	}
// 	tmp.Name = p.Name
// 	tmp.Controls = p.Controls
// 	return json.MarshalIndent(&tmp, "", "  ")
// }
