package registry

import (
	"context"
	"errors"
	"log"
	"net/url"

	"github.com/manifoldco/torus-cli/envelope"
	"github.com/manifoldco/torus-cli/identity"
	"github.com/manifoldco/torus-cli/primitive"
)

// MembershipsClient represents the `/memberships` registry
// endpoint, used for accessing the relationship between users,
// organization, and teams.
type MembershipsClient struct {
	client RoundTripper
}

// List returns all memberships for a given organization, team, or user/machine
func (m *MembershipsClient) List(ctx context.Context, orgID *identity.ID,
	teamID *identity.ID, ownerID *identity.ID) ([]envelope.Membership, error) {

	query := &url.Values{}
	if orgID != nil {
		query.Set("org_id", orgID.String())
	}
	if teamID != nil {
		query.Set("team_id", teamID.String())
	}
	if ownerID != nil {
		query.Set("owner_id", ownerID.String())
	}

	req, err := m.client.NewRequest("GET", "/memberships", query, nil)
	if err != nil {
		log.Printf("could not build GET /memberships request: %s", err)
		return nil, err
	}

	memberships := []envelope.Membership{}
	_, err = m.client.Do(ctx, req, &memberships)
	if err != nil {
		log.Printf("could not perform GET /memberships: %s", err)
		return nil, err
	}

	return memberships, nil
}

// Create requests addition of a user to a team
func (m *MembershipsClient) Create(ctx context.Context, userID, orgID, teamID *identity.ID) error {
	if orgID == nil {
		return errors.New("invalid org")
	}
	if userID == nil {
		return errors.New("invalid user")
	}
	if teamID == nil {
		return errors.New("invalid team")
	}

	membershipBody := primitive.Membership{
		OwnerID: userID,
		OrgID:   orgID,
		TeamID:  teamID,
	}

	ID, err := identity.NewMutable(&membershipBody)
	if err != nil {
		return err
	}

	membership := envelope.Membership{
		ID:      &ID,
		Version: 1,
		Body:    &membershipBody,
	}

	req, err := m.client.NewRequest("POST", "/memberships", nil, membership)
	if err != nil {
		return err
	}

	_, err = m.client.Do(ctx, req, nil)
	return err
}

// Delete requests deletion of a specific membership row by ID
func (m *MembershipsClient) Delete(ctx context.Context, membership *identity.ID) error {
	req, err := m.client.NewRequest("DELETE", "/memberships/"+membership.String(), nil, nil)
	if err != nil {
		return err
	}

	_, err = m.client.Do(ctx, req, nil)
	return err
}
