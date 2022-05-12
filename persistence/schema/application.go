package schema

import "fmt"

const (
	// We use `-1` with semantic variable names to indicate non-existence or non-validity
	// in various contexts to avoid the usage of nullability in columns and for performance
	// optimization purposes.
	DefaultUnstakingHeight = -1 // TODO(team): Move this into a shared file?
	DefaultEndHeight       = -1 // TODO(team): Move this into a shared file?

	// TODO (team) look into address being a "computed" field
	// DISCUSS(drewsky): How do we handle historical queries here? E.g. get staked chains at some specific height?
	AppTableName   = "app"
	AppTableSchema = `(
			address    	     TEXT NOT NULL,
			public_key 		 TEXT NOT NULL,
			staked_tokens    TEXT NOT NULL,
			max_relays		 TEXT NOT NULL,
			output_address   TEXT NOT NULL,
			paused_height    BIGINT NOT NULL default -1,
			unstaking_height BIGINT NOT NULL default -1,
			end_height       BIGINT NOT NULL default -1,

			/* DISCUSS(drewsky): We can't do ON CONFLICT multiple constraints, so what should we do here? */
			/* CONSTRAINT app_paused_height UNIQUE (address, paused_height), */
			/* CONSTRAINT app_paused_height UNIQUE (address, unstaking_height), */
			CONSTRAINT app_end_height UNIQUE (address, end_height)
		)`

	AppChainsTableName   = "app_chains"
	AppChainsTableSchema = `(
			address      TEXT NOT NULL,
			chain_id     CHAR(4) NOT NULL,
			end_height   BIGINT NOT NULL default -1,

			CONSTRAINT app_chain_end_height UNIQUE (address, chain_id, end_height)
		)`
)

func AppQuery(address string) string {
	return fmt.Sprintf(`SELECT * FROM %s WHERE address='%s' AND end_height=%d`, AppTableName, address, DefaultEndHeight)
}

func AppChainsQuery(address string) string {
	return fmt.Sprintf(`SELECT * FROM %s WHERE address='%s' AND end_height=%d`, AppChainsTableName, address, DefaultEndHeight)
}

// DISCUSS(drewsky): Do we not want to filter by DefaultEndHeight=-1 here?
func AppExistsQuery(address string) string {
	return fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE address='%s')`, AppTableName, address)
}

// DISCUSS(drewsky): Do we not want to filter by `unstaking_height >= unstakingHeight AND end_height=DefaultEndHeight=-1" here?
func AppReadyToUnstakeQuery(unstakingHeight int64) string {
	return fmt.Sprintf(`SELECT address, staked_tokens, output_address FROM %s WHERE unstaking_height=%d`, AppTableName, unstakingHeight)
}

func AppOutputAddressQuery(operatorAddress string, height int64) string {
	return fmt.Sprintf(`SELECT output_address FROM %s WHERE address='%s' AND end_height=%d`,
		AppTableName, operatorAddress, DefaultEndHeight)
}

// DISCUSS(team): if current_height == unstaking_height - is the actor unstaking or unstaked (i.e. did we process the block yet)?
func AppUnstakingHeightQuery(address string, height int64) string {
	return fmt.Sprintf(`SELECT unstaking_height FROM %s WHERE address='%s' AND end_height=%d`,
		AppTableName, address, DefaultEndHeight)
}

func AppPauseHeightQuery(address string, height int64) string {
	return fmt.Sprintf(`SELECT paused_height FROM %s WHERE address='%s' AND end_height=%d`,
		AppTableName, address, DefaultEndHeight)
}

func InsertAppQuery(address, publicKey, stakedTokens, maxRelays, outputAddress string, pausedHeight, unstakingHeight int64, chains []string) string {
	// insert the app
	insertIntoAppTable := fmt.Sprintf(
		`WITH
			ins1 AS (INSERT INTO %s(address, public_key, staked_tokens, max_relays, output_address, paused_height, unstaking_height, end_height)
			VALUES('%s','%s','%s','%s','%s',%d,%d,%d)
			RETURNING address)`,
		AppTableName, address, publicKey, stakedTokens, maxRelays, outputAddress, pausedHeight, unstakingHeight, DefaultEndHeight)

	// Insert each app chain
	maxIndex := len(chains) - 1
	insertIntoAppTable += "\nINSERT INTO app_chains (address, chain_id, end_height) VALUES"
	for i, chain := range chains {
		insertIntoAppTable += fmt.Sprintf("\n((SELECT address FROM ins1), '%s', %d)", chain, DefaultEndHeight)
		if i < maxIndex {
			insertIntoAppTable += ","
		}
	}
	return insertIntoAppTable
}

func NullifyAppQuery(address string, height int64) string {
	return fmt.Sprintf(`UPDATE %s SET end_height=%d WHERE address='%s' AND end_height=%d`,
		AppTableName, height, address, DefaultEndHeight)
}

func NullifyAppChainsQuery(address string, height int64) string {
	return fmt.Sprintf(`UPDATE %s SET end_height=%d WHERE address='%s' AND end_height=%d`,
		AppChainsTableName, height, address, DefaultEndHeight)
}

func UpdateAppQuery(address, stakedTokens, maxRelays string, height int64) string {
	return fmt.Sprintf(
		`INSERT INTO %s(address, public_key, staked_tokens, max_relays, output_address, paused_height, unstaking_height, end_height)
			(
				SELECT address, public_key, '%s', '%s', output_address, paused_height, unstaking_height, %d
				FROM %s WHERE address='%s' AND (end_height=%d OR end_height=%d)
			)
		ON CONFLICT ON CONSTRAINT app_end_height
			DO UPDATE SET staked_tokens=EXCLUDED.staked_tokens, max_relays=EXCLUDED.max_relays, end_height=EXCLUDED.end_height`,
		AppTableName,
		stakedTokens, maxRelays, DefaultEndHeight,
		AppTableName, address, height, DefaultEndHeight)
}

func UpdateAppUnstakingHeightQuery(address string, unstakingHeight, height int64) string {
	return fmt.Sprintf(`
		INSERT INTO %s(address, public_key, staked_tokens, max_relays, output_address, paused_height, unstaking_height, end_height)
			(
				SELECT address, public_key, staked_tokens, max_relays, output_address, paused_height, %d, %d
				FROM %s WHERE address='%s' AND (end_height=%d OR end_height=%d)
			)
		ON CONFLICT ON CONSTRAINT app_end_height
			DO UPDATE SET unstaking_height=EXCLUDED.unstaking_height, end_height=EXCLUDED.end_height`,
		AppTableName,
		unstakingHeight, DefaultEndHeight,
		AppTableName, address, height, DefaultEndHeight)

}

func UpdateAppPausedHeightQuery(address string, pausedHeight, height int64) string {
	return fmt.Sprintf(`
		INSERT INTO %s(address, public_key, staked_tokens, max_relays, output_address, paused_height, unstaking_height, end_height)
			(
				SELECT address, public_key, staked_tokens, max_relays, output_address, %d, unstaking_height, %d
				FROM %s WHERE address='%s' AND (end_height=%d OR end_height=%d)
			)
		ON CONFLICT ON CONSTRAINT app_end_height
			DO UPDATE SET paused_height=EXCLUDED.paused_height, end_height=EXCLUDED.end_height`,
		AppTableName,
		pausedHeight, DefaultEndHeight,
		AppTableName, address, height, DefaultEndHeight)
}

func UpdateAppsPausedBefore(pauseBeforeHeight, unstakingHeight, currentHeight int64) string {
	return fmt.Sprintf(`
		INSERT INTO %s(address, public_key, staked_tokens, max_relays, output_address, paused_height, unstaking_height, end_height)
			(
				SELECT address, public_key, staked_tokens, max_relays, output_address, paused_height, %d, %d
				FROM %s WHERE paused_height<%d AND paused_height>=0 AND (end_height=%d OR end_height=%d)
			)
		ON CONFLICT ON CONSTRAINT app_end_height
			DO UPDATE SET unstaking_height=EXCLUDED.unstaking_height, end_height=EXCLUDED.end_height`,
		AppTableName,
		unstakingHeight, DefaultEndHeight,
		AppTableName, pauseBeforeHeight, currentHeight, DefaultEndHeight)
}

func NullifyAppsPausedBeforeQuery(pausedBeforeHeight, height int64) string {
	return fmt.Sprintf(`
		UPDATE %s SET end_height=%d
		WHERE paused_height<%d AND paused_height>=0 AND end_height=%d`,
		AppTableName, height, pausedBeforeHeight, DefaultEndHeight)
}

func UpdateAppChainsQuery(address string, chains []string, height int64) string {
	insert := fmt.Sprintf("INSERT INTO %s (address, chain_id, end_height) VALUES", AppChainsTableName)
	maxIndex := len(chains) - 1
	for i, chain := range chains {
		insert += fmt.Sprintf("\n('%s', '%s', %d)", address, chain, DefaultEndHeight)
		if i < maxIndex {
			insert += ","
		}
	}
	return insert
}

func ClearAllAppsQuery() string {
	return fmt.Sprintf(`DELETE FROM %s`, AppTableName)
}

func ClearAllAppChainsQuery() string {
	return fmt.Sprintf(`DELETE FROM %s`, AppChainsTableName)
}
