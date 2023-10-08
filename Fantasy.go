package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "time"
	"errors"

    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
)

const (
    // MySQL Database configuration
    dbHost     = "localhost"
    dbPort     = 3306
    dbUser     = "your_username"
    dbPassword = "your_password"
    dbName     = "your_database"
)

type Contest struct {
    ID           int
    Name         string
    Prize        float64
    TotalSlots   int
    RemainingSlots int
    StartDate    time.Time
    EndDate      time.Time
    Status       string
    ActiveDate   time.Time
    CreatedAt    time.Time
}

type Team struct {
    ID          int
    Name        string
    DisplayName string
    CreatedAt   time.Time
}

func main() {
    // Establish a database connection
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    r := gin.Default()

    // Define your API routes here

    r.Run(":8080")
}

// CRUD operations for contests

// Create a new contest
func createContest(db *sql.DB, name string, prize float64, totalSlots int, startDate, endDate time.Time) error {
	// Create a transaction to ensure consistency
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	// Insert a new contest record into the database
	_, err = tx.Exec(
		"INSERT INTO contest (name, prize, total_slots, remaining_slots, start_date, end_date, status, active_date, created_at) VALUES (?, ?, ?, ?, ?, ?, 'active', ?, NOW())",
		name, prize, totalSlots, totalSlots, startDate, endDate, startDate,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}



// Get a contest by ID
func getContest(db *sql.DB, contestID int) (*Contest, error) {
	// Query the database to fetch the contest by ID
	row := db.QueryRow("SELECT id, name, prize, total_slots, remaining_slots, start_date, end_date, status, active_date, created_at FROM contest WHERE id = ?", contestID)

	var contest Contest
	err := row.Scan(
		&contest.ID,
		&contest.Name,
		&contest.Prize,
		&contest.TotalSlots,
		&contest.RemainingSlots,
		&contest.StartDate,
		&contest.EndDate,
		&contest.Status,
		&contest.ActiveDate,
		&contest.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// No contest found with the given ID
			return nil, errors.New("contest not found")
		}
		// Other database-related error
		return nil, err
	}

	return &contest, nil
}



// Update a contest's slot
func updateContestSlot(db *sql.DB, contestID int, newSlot int) error {
	// Start a transaction to ensure consistency
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	// Check if the contest with the given ID exists
	var currentSlots int
	err = tx.QueryRow("SELECT remaining_slots FROM contest WHERE id = ?", contestID).Scan(&currentSlots)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Calculate the new remaining slots
	newRemainingSlots := currentSlots + (newSlot - currentSlots)

	// Update the contest's remaining slots
	_, err = tx.Exec("UPDATE contest SET remaining_slots = ? WHERE id = ?", newRemainingSlots, contestID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Delete a contest by ID
func deleteContest(db *sql.DB, contestID int) error {
	// Start a transaction to ensure consistency
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	// Delete the contest with the given ID from the database
	_, err = tx.Exec("DELETE FROM contest WHERE id = ?", contestID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}


// CRUD operations for teams

// Create a new team
func createTeam(db *sql.DB, name, displayName string) error {
	// Start a transaction to ensure consistency
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	// Insert a new team record into the database
	_, err = tx.Exec("INSERT INTO team (name, displayname, created_at) VALUES (?, ?, NOW())", name, displayName)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}



// Get a team by ID
func getTeam(db *sql.DB, teamID int) (*Team, error) {
    // Query the database to fetch the team by ID
    row := db.QueryRow("SELECT id, name, displayname, created_at FROM team WHERE id = ?", teamID)

    var team Team
    err := row.Scan(
        &team.ID,
        &team.Name,
        &team.DisplayName,
        &team.CreatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            // No team found with the given ID
            return nil, errors.New("team not found")
        }
        // Other database-related error
        return nil, err
    }

    return &team, nil
}


func setupRoutes(r *gin.Engine) {
		// Route to create a new team
		r.POST("/teams", func(c *gin.Context) {
			// Parse the request body to get the team data
			var team Team // Assuming you have a Team struct defined
	
			if err := c.ShouldBindJSON(&team); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
	
			// Validate and create the team in the database
			if err := createTeam(db, team.Name, team.DisplayName); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
				return
			}
	
			c.JSON(http.StatusCreated, gin.H{"message": "Team created successfully"})
		})
	}
	
	
	func setupRoutes(r *gin.Engine) {
		// Route to fetch all teams
		r.GET("/teams", func(c *gin.Context) {
			// logic to fetch all teams from the database
	
			// Assuming you have a function getTeams(db *sql.DB) that fetches all teams
			teams, err := getTeams(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teams"})
				return
			}
	
			c.JSON(http.StatusOK, teams)
		})
	}
	
	
	func setupRoutes(r *gin.Engine, db *sql.DB) {
		// Route to fetch a particular team by ID
		r.GET("/teams/:id", func(c *gin.Context) {
			// Get the team ID from the URL parameter
			teamIDStr := c.Param("id")
	
			// Convert the teamIDStr to an integer
			teamID, err := strconv.Atoi(teamIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
				return
			}
	
			// logic to fetch a team by ID from the database
			team, err := getTeamByID(db, teamID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team"})
				return
			}
	
			c.JSON(http.StatusOK, team)
		})
	}
	
	// Define your Team struct here if not already defined
	type Team struct {
		ID          int
		Name        string
		DisplayName string
		CreatedAt   time.Time
	}
	
	// getTeamByID function that retrieves a team by its ID
	func getTeamByID(db *sql.DB, teamID int) (*Team, error) {
		var team Team
		err := db.QueryRow("SELECT id, name, displayname, created_at FROM team WHERE id = ?", teamID).
			Scan(&team.ID, &team.Name, &team.DisplayName, &team.CreatedAt)
	
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("team not found")
			}
			return nil, err
		}
	
		return &team, nil
	}
	
	
	func setupRoutes(r *gin.Engine, db *sql.DB) {
		// Route to create a new contest
		r.POST("/contests", func(c *gin.Context) {
			// Parse the request body to get the contest data
			var contest Contest // Assuming you have a Contest struct defined
	
			if err := c.ShouldBindJSON(&contest); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
	
			// Set default values for status and active_date
			contest.Status = "active"
			contest.ActiveDate = time.Now()
	
			// Validate and create the contest in the database
			if err := createContest(db, contest); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contest"})
				return
			}
	
			c.JSON(http.StatusCreated, gin.H{"message": "Contest created successfully"})
		})
	}
	
	// Define your Contest struct here if not already defined
	type Contest struct {
		Name         string    `json:"name"`
		Prize        float64   `json:"prize"`
		TotalSlots   int       `json:"total_slots"`
		RemainingSlots int     `json:"remaining_slots"`
		StartDate    time.Time `json:"start_date"`
		EndDate      time.Time `json:"end_date"`
		Status       string    `json:"status"`
		ActiveDate   time.Time `json:"active_date"`
	}
	
	// createContest function that inserts a new contest into the database
	func createContest(db *sql.DB, contest Contest) error {
		// Start a transaction to ensure consistency
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
			}
		}()
	
		// Insert a new contest record into the database
		_, err = tx.Exec(
			"INSERT INTO contest (name, prize, total_slots, remaining_slots, start_date, end_date, status, active_date, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())",
			contest.Name, contest.Prize, contest.TotalSlots, contest.TotalSlots, contest.StartDate, contest.EndDate, contest.Status, contest.ActiveDate,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	
		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			return err
		}
	
		return nil
	}
	

			func setupRoutes(r *gin.Engine, db *sql.DB) {
				// Route to enter a contest
				r.POST("/contests/enter", func(c *gin.Context) {
					// Parse the request body to get the user's entry data
					var entry ContestEntry
			
					if err := c.ShouldBindJSON(&entry); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
			
					// Assuming you have a function enterContest(db *sql.DB, entry ContestEntry) that handles contest entry
					if err := enterContest(db, entry); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enter contest"})
						return
					}
			
					c.JSON(http.StatusCreated, gin.H{"message": "Entered contest successfully"})
				})
			}
			
			// Define your ContestEntry struct here if not already defined
			type ContestEntry struct {
				ContestID int `json:"contest_id"`
				UserID    int `json:"user_id"`
			}
			
			// enterContest function that handles contest entry
			func enterContest(db *sql.DB, entry ContestEntry) error {
				// Start a transaction to ensure consistency
				tx, err := db.Begin()
				if err != nil {
					return err
				}
				defer func() {
					if p := recover(); p != nil {
						tx.Rollback()
					}
				}()
			
				// 1. Check User Eligibility (Define your eligibility criteria)
				if !isUserEligible(db, entry.UserID, entry.ContestID) {
					tx.Rollback()
					return errors.New("User is not eligible to enter this contest")
				}
			
				// 2. Check Remaining Slots
				remainingSlots, err := getRemainingSlots(db, entry.ContestID)
				if err != nil {
					tx.Rollback()
					return err
				}
			
				if remainingSlots <= 0 {
					tx.Rollback()
					return errors.New("No remaining slots available in the contest")
				}
			
				// 3. Update Contest Slots
				_, err = tx.Exec("UPDATE contest SET remaining_slots = remaining_slots - 1 WHERE id = ?", entry.ContestID)
				if err != nil {
					tx.Rollback()
					return err
				}
			
				// 4. Insert a record in the user-contest relationship table
				_, err = tx.Exec("INSERT INTO user_contest (user_id, contest_id) VALUES (?, ?)", entry.UserID, entry.ContestID)
				if err != nil {
					tx.Rollback()
					return err
				}
			
				// Commit the transaction
				err = tx.Commit()
				if err != nil {
					return err
				}
			
				return nil
			}
			
			// Define your eligibility criteria function (isUserEligible) and remaining slots retrieval function (getRemainingSlots) here.
			// Define your eligibility criteria function (isUserEligible) here.
func isUserEligible(db *sql.DB, userID int, contestID int) bool {
    // Example eligibility criteria: Check if the user is at least 18 years old
    minAge := 18

    var userAge int
    err := db.QueryRow("SELECT age FROM users WHERE id = ?", userID).Scan(&userAge)
    if err != nil {
        // Handle the error (e.g., user not found, database error)
        return false
    }

    return userAge >= minAge
}

// Define your remaining slots retrieval function (getRemainingSlots) here.
func getRemainingSlots(db *sql.DB, contestID int) (int, error) {
    // Query the remaining slots from the "contest" table
    var remainingSlots int
    err := db.QueryRow("SELECT remaining_slots FROM contest WHERE id = ?", contestID).Scan(&remainingSlots)
    if err != nil {
        // Handle the error (e.g., contest not found, database error)
        return 0, err
    }

    return remainingSlots, nil
}

	
	func setupRoutes(r *gin.Engine, db *sql.DB) {
		// Route to change the selected contest for a user
		r.PUT("/contests/change/:userID", func(c *gin.Context) {
			// Parse the request body to get the new contest selection
			var contestChange ContestChange // Assuming you have a ContestChange struct defined
	
			if err := c.ShouldBindJSON(&contestChange); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
	
			// Get the user ID from the URL parameter
			userIDStr := c.Param("userID")
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}
	
			// changeSelectedContest function that handles changing the selected contest for a user
		func changeSelectedContest(db *sql.DB, userID int, newContestID int) error {
    	// Start a database transaction to ensure data consistency
    	tx, err := db.Begin()
    	if err != nil {
	        return err
    	}
    	defer func() {
        	if p := recover(); p != nil {
            	tx.Rollback()
        	}
    	}()

    // Check if the new contest is valid and exists
    if !isContestValid(db, newContestID) {
        tx.Rollback()
        return errors.New("Invalid or non-existent contest selected")
    }

    // Update the user's selected contest in the database
    _, err = tx.Exec("UPDATE users SET selected_contest_id = ? WHERE id = ?", newContestID, userID)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Commit the transaction
    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

// check if the new contest is valid and exists
func isContestValid(db *sql.DB, contestID int) bool {
    var valid bool
    err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM contest WHERE id = ?)", contestID).Scan(&valid)
    if err != nil {
        return false
    }
    return valid
}

	
			// Assuming you have a function changeSelectedContest(db *sql.DB, userID int, newContestID int) that handles the contest change
			if err := changeSelectedContest(db, userID, contestChange.NewContestID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change the selected contest"})
				return
			}
	
			c.JSON(http.StatusOK, gin.H{"message": "Selected contest changed successfully"})
		})
	}
	
	// Define your ContestChange struct here if not already defined
	type ContestChange struct {
		NewContestID int `json:"new_contest_id"`
	}
	
	// changeSelectedContest function that handles changing the selected contest for a user
	func changeSelectedContest(db *sql.DB, userID int, newContestID int) error {
		// This may involve updating the user's contest selection in the database
		// logic to change the selected contest for a user
func changeSelectedContest(db *sql.DB, userID int, newContestID int) error {
    // Start a transaction to ensure data consistency
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
        }
    }()

    // Check if the new contest (newContestID) exists in the database
    var contestExists bool
    err = tx.QueryRow("SELECT EXISTS (SELECT 1 FROM contest WHERE id = ?)", newContestID).Scan(&contestExists)
    if err != nil {
        tx.Rollback()
        return err
    }

    if !contestExists {
        tx.Rollback()
        return errors.New("The specified contest does not exist")
    }

    // Update the user's selected contest in the users table
    _, err = tx.Exec("UPDATE users SET selected_contest_id = ? WHERE id = ?", newContestID, userID)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Commit the transaction
    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

	
		_, err := db.Exec("UPDATE user_contest SET contest_id = ? WHERE user_id = ?", newContestID, userID)
		if err != nil {
			return err
		}
	
		return nil
	}
	
	
	func setupRoutes(r *gin.Engine, db *sql.DB) {
		// Route to leave a contest
		r.DELETE("/contests/leave/:userID", func(c *gin.Context) {
			// Get the user ID from the URL parameter
			userIDStr := c.Param("userID")
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}
	
			// logic to allow users to leave a contest
func leaveContest(db *sql.DB, userID int) error {
    // Start a database transaction to ensure data consistency
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
        }
    }()

    // logic to allow users to leave a contest here
    // Option 1: Update the user's selected contest to null (assuming NULL is used to indicate no selection)
    _, err = tx.Exec("UPDATE users SET selected_contest_id = NULL WHERE id = ?", userID)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Option 2: Delete the user's participation record in the user-contest relationship table
    _, err = tx.Exec("DELETE FROM user_contest WHERE user_id = ?", userID)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Commit the transaction
    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

	
			// Assuming you have a function leaveContest(db *sql.DB, userID int) that handles contest leaving
			if err := leaveContest(db, userID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave the contest"})
				return
			}
	
			c.JSON(http.StatusOK, gin.H{"message": "Left contest successfully"})
		})
	}
	
	// leaveContest function that handles contest leaving
	func leaveContest(db *sql.DB, userID int) error {
		// Start a database transaction to ensure data consistency
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
			}
		}()
	
		// This may involve updating the user's contest selection or removing the user from the contest participation records
		// logic to allow users to leave a contest
func leaveContest(db *sql.DB, userID int, contestID int) error {
    // Start a database transaction to ensure data consistency
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
        }
    }()

    // Check if the user is currently participating in the contest
    var isParticipating bool
    err = tx.QueryRow("SELECT EXISTS (SELECT 1 FROM user_contest WHERE user_id = ? AND contest_id = ?)", userID, contestID).Scan(&isParticipating)
    if err != nil {
        tx.Rollback()
        return err
    }

    if !isParticipating {
        tx.Rollback()
        return errors.New("User is not participating in the contest")
    }

    // logic to allow users to leave a contest here

    // Option 1: Update the user's selected contest to null (assuming NULL is used to indicate no selection)
    _, err = tx.Exec("UPDATE users SET selected_contest_id = NULL WHERE id = ?", userID)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Option 2: Delete the user's participation record in the user-contest relationship table
    _, err = tx.Exec("DELETE FROM user_contest WHERE user_id = ? AND contest_id = ?", userID, contestID)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Commit the transaction
    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}

	
		// For example, you can update the user's selected contest to null (assuming NULL is used to indicate no selection):
		_, err = tx.Exec("UPDATE users SET selected_contest_id = NULL WHERE id = ?", userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	
		// Or you can delete the user's participation record in the user-contest relationship table:
		_, err = tx.Exec("DELETE FROM user_contest WHERE user_id = ?", userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	
		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			return err
		}
	
		return nil
	}
	
