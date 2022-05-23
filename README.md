# go-pingdom #

[![GoDoc](https://godoc.org/github.com/algolia/go-pingdom/pingdom?status.svg)](https://godoc.org/github.com/algolia/go-pingdom/pingdom)

go-pingdom is a Go client library for the Pingdom API.

This currently supports working with HTTP, ping checks, and TCP checks.

**Important**: The current version of this library only supports the Pingdom 3.1 API.  If you are still using the deprecated Pingdom 2.1 API please pin your dependencies to tag v1.1.0 of this library.

## Usage ##

### Pingdom Client ###

Construct a new Pingdom client:

```go
client, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
    APIToken: "pingdom_api_token",
})
```

Using a Pingdom client, you can access supported services.

You can override the timeout or other parameters by passing a custom http client:
```go
client, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
    APIToken: "pingdom_api_token",
    HTTPClient: &http.Client{
        Timeout: time.Second * 10,
    },
})
```

The `APIToken` can also implicitly be provided by setting the environment variable `PINGDOM_API_TOKEN`:

```bash
export PINGDOM_API_TOKEN=pingdom_api_token
./your_application
```


### Pindom Extension Client ###

Construct a new Pingdom extension client:

```go
client_ext, err := pingdomext.NewClientWithConfig(pingdomext.ClientConfig{
    Username: "test_user",
    Password: "test_pwd",
    OrgId: "test_org"
    HTTPClient: &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
    },
})
```

Using a Pingdom extention client, you can access supported services, like integration service.

You must override the CheckRedirect since there have multiple redirect while get the jwt token for access api. 

The `Username`,`Password` and `OrgID`can also implicitly be provided by setting the environment variable `SOLARWINDS_USER` , `SOLARWINDS_PASSWD` and `SOLARWINDS_ORG_ID`:

```bash
export SOLARWINDS_USER=test_user
export SOLARWINDS_PASSWD=test_pwd
export SOLARWINDS_ORG_ID=test_org
./your_application
```

The `Username` and `Password` is required, the `OrgID` is optional. If the `OrgID` is not provide, your default organization will be used.

### Solarwinds Client ###

Construct a new Solarwinds client:

```go
solarwindsClient, err := solarwinds.NewClient(solarwinds.ClientConfig{
    Username: "solarwinds web portal login username",
    Password: "solarwinds web portal login password"
})
```

Using a Solarwinds client, you can access supported services.

### CheckService ###

This service manages pingdom Checks which are represented by the `Check` struct.
When creating or updating Checks you must specify at a minimum the `Name`, `Hostname`
and `Resolution`.  Other fields are optional but if not set will be given the zero
values for the underlying type.

More information on Checks from Pingdom: https://www.pingdom.com/features/api/documentation/#ResourceChecks

Get a list of all checks:

```go
checks, err := client.Checks.List()
fmt.Println("Checks:", checks) // [{ID Name} ...]
```

Create a new HTTP check:

```go
newCheck := pingdom.HttpCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Create a new Ping check:
```go
newCheck := pingdom.PingCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Create a new TCP check:
```go
newCheck := pingdom.TCPCheck{Name: "Test Check", Hostname: "example.com", Port: 25, StringToSend: "HELO foo.com", StringToExpect: "250 mail.test.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Create a new DNS check:
```go
newCheck := pingdom.DNSCheck{
    Name: "fake check",
    Hostname: "example.com",
    ExpectedIP: "192.168.1.1",
    NameServer: "8.8.8.8",
}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Get details for a specific check:

```go
checkDetails, err := client.Checks.Read(12345)
```

For checks with detailed information, check the specific details in
the field `Type` (e.g. `checkDetails.Type.HTTP`).

Update a check:

```go
updatedCheck := pingdom.HttpCheck{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
msg, err := client.Checks.Update(12345, &updatedCheck)
```

Delete a check:

```go
msg, err := client.Checks.Delete(12345)
```

Create a check with basic alert notification to a user.

```go
newCheck := pingdom.HttpCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5, SendNotificationWhenDown: 2, UserIds []int{12345}}
checkResponse, err := client.Checks.Create(&newCheck)
```

### MaintenanceService ###

This service manages pingdom Maintenances which are represented by the `Maintenance` struct.
When creating or updating Maintenances you must specify at a minimum the `Description`, `From`
and `To`.  Other fields are optional but if not set will be given the zero
values for the underlying type.

More information on Maintenances from Pingdom: https://www.pingdom.com/resources/api/2.1#ResourceMaintenance

Get a list of all maintenances:

```go
maintenances, err := client.Maintenances.List()
fmt.Println("Maintenances:", maintenances) // [{ID Description} ...]
```

Create a new Maintenance Window:

```go
m := pingdom.MaintenanceWindow{
    Description: "My Maintenance",
    From:        1,
    To:          1234567899,
}
maintenance, err := client.Maintenances.Create(&m)
fmt.Println("Created MaintenanceWindow:", maintenance) // {ID Description}
```

Get details for a specific maintenance:

```go
maintenance, err := client.Maintenances.Read(12345)
```

Update a maintenance: (Please note, that based on experience, you are allowed to modify only `Description`, `EffectiveTo` and `To`)

```go
updatedMaintenance := pingdom.MaintenanceWindow{
    Description: "My Maintenance",
    To:          1234567999,
}
msg, err := client.Maintenances.Update(12345, &updatedMaintenance)
```

Delete a maintenance:

Note: that only future maintenance window can be deleted. This means that both `To` and `From` should be in future.

```go
msg, err := client.Maintenances.Delete(12345)
```

After contacting Pingdom, the better approach would be to use update function and setting `To` and `EffectiveTo` to current time

```go
maintenance, _ := client.Maintenances.Read(12345)

m := pingdom.MaintenanceWindow{
    Description: maintenance.Description,
    From:        maintenance.From,
    To:          1,
    EffectiveTo: 1,
}

maintenanceUpdate, err := client.Maintenances.Update(12345, &m)
```

### OccurrenceService ###

This service manages pingdom Maintenance Occurrences which are represented by the `Occurrence` struct.
It is not possible to create occurrences directly, instead they are created automatically as specified by
maintenances. Only `From` and `To` fields can be updated when updating an occurrence.

More information on Occurrences from Pingdom: https://docs.pingdom.com/api/#tag/Maintenance-occurrences

Get a list of all occurrences:

```go
occurrences, err := client.Occurrences.List(ListOccurrenceQuery{})
fmt.Println("Occurrences:", occurrences) // [{ID Description} ...]
```

Get details for a specific occurrence:

```go
occurrence, err := client.Occurrences.Read(12345)
```

Update an occurrence: (Please note, that based on experience, you are allowed to modify only `From` and `To`)

Note: that only future maintenance occurences can be updated.

```go
update := pingdom.Occurrence{
    From:        1,
    To:          1234567999,
}
msg, err := client.Occurrences.Update(12345, update)
```

Delete an Occurrence:

Note: that only future maintenance occurrence can be deleted. 

```go
msg, err := client.Maintenances.Delete(12345)
```

Delete multiple Occurrences in one go:

```go
msg, err := client.Maintenances.Delete([]int64{1, 2, 3, 4, 5})
```

### ProbeService ###

This service gets pingdom Probes which are represented by the `Probes` struct.

More information on Probes from Pingdom: https://www.pingdom.com/resources/api/2.1#ResourceProbes
Several parameters are supported for filtering output. Please see them in Pingdom API documentation.

**NOTE:** Official documentation does not specify that `region` is returned for every probe entry, but it does and you can use it.

Get a list of all probes:

```go
params := make(map[string]string)

probes, err := client.Probes.List(params)
fmt.Println("Probes:", probes) // [{ID Name} ...]

for _, probe := range probes {
    fmt.Println("Probe region:", probe.Region)  // Probe region: EU
}
```

### TeamService ###

This service manages pingdom Teams which are represented by the `Team` struct.
When creating or updating Teams you must specify the `Name` and `MemberIDs`,
though `MemberIDs` may be an empty slice.
More information on Teams from Pingdom: https://docs.pingdom.com/api/#tag/Teams

Get a list of all teams:

```go
teams, err := client.Teams.List()
fmt.Println("Teams:", teams) // [{ID Name MemberIDs} ...]
```

Create a new Team:

```go
t := pingdom.TeamData{
    Name: "Team",
    MemberIDs: []int{},
}
team, err := client.Teams.Create(&t)
fmt.Println("Created Team:", team) // {ID Name MemberIDs}
```

Get details for a specific team:

```go
team, err := client.Teams.Read(12345)
```

Update a team:

```go
modifyTeam := pingdom.TeamData{
    Name:    "New Name"
    MemberIDs: []int{123, 678},
}
team, err := client.Teams.Update(12345, &modifyTeam)
```

Delete a team:

```go
team, err := client.Teams.Delete(12345)
```

### ContactService ###

This service manages users and their contact information which is represented by the `Contact` struct.
More information from Pingdom: https://docs.pingdom.com/api/#tag/Contacts

Get all contact info:

```go
contacts, err := client.Contacts.List()
fmt.Println(contacts)
```

Create a new contact:

```go
contact := Contact{
    Name: "John Doe",
    Paused: false,
    NotificationTargets: NotificationTargets{
        SMS: []SMSNotificationTarget{
            {
                Number: "5555555555",
                CountryCode: "1",
                Provider: "Verizon",
            }
        }
    }
}
contactId, err := client.Contacts.Create(contact)
fmt.Println("New Contact ID: ", contactId.Id)
```

Update a contact

```go
contactId := 1234

contact := Contact{
    Name : "John Doe",
    Paused : false,
    NotificationTargets: NotificationTargets{
        SMS: []SMSNotificationTarget{
            {
                Number: "5555555555",
                CountryCode: "1",
                Provider: "T-Mobile",
            }
        }
    }
}
result, err := client.Contacts.Update(contactId, contact)
fmt.Println(result.Message)
```

Delete a contact

```go
contactId := 1234

result, err := client.Contacts.Delete(contactId)
fmt.Println(result.Message)
```

### TMS Checks Service ###

This service manages pingdom TMS Checks which are represented by the `TMS Check` struct.
More information from Pingdom: https://docs.pingdom.com/api/#tag/TMS-Checks


Get a list of all TMS Checks:

```go
tmsChecks, err := client.TMSCheck.List()
fmt.Println("TMS Checks:", tmsChecks) 
```

Create a new TMS Check:

```go
tmsCheck := pingdom.TMSCheck{
		Name: "wlwu-test-111",
		Steps: []pingdom.TMSCheckStep{
			{
				Args: map[string]string{
					"url": "www.google.com",
				},
				Fn: "go_to",
			},
		},
	}

createMsg, err := client.TMSCheck.Create(&tmsCheck)
tmsCheckID := createMsg.ID
```

Get details for a specific TMS Check:

```go
tmsCheckDetail, err := client.TMSCheck.Read(12345)
```


Update a TMS Check:

```go
tmsCheck := pingdom.TMSCheck{
		Name: "wlwu-test-222",
		Steps: []pingdom.TMSCheckStep{
			{
				Args: map[string]string{
					"url": "www.google.com",
				},
				Fn: "go_to",
			},
		},
	}
updateMsg, err := client.TMSCheck.Update(12345, &tmsCheck)
```

Delete a TMS Check:

```go
delMsg, err := client.TMSCheck.Delete(12345)
```





### IntegrationService ###

This service manages pingdom Integrations which are represented by the `Integration` struct. Now only support manages the WebHook Integrations.
When creating or updating Integrations you must specify the `Active`, `ProviderID` and `WebHookData`.  


Get a list of all integrations:

```go
integrations, err := client_ext.Integrations.List()
fmt.Println("Integrations:", integrations) 
```

Create a new WebHook Integration:

```go
newIntegration := pingdomext.WebHookIntegration{
	Active:     false,
	ProviderID: 2,
	UserData: &pingdomext.WebHookData{
		Name: "tets-1",
		URL:  "http://www.example.com",
	},
}
integrationStatus, err := client_ext.Integrations.Create(&newIntegration)
fmt.Println("Created integration:", integrationStatus) 
```

Get details for a specific integration:

```go
integrationDetail, err := client_ext.Integrations.Read(12345)
```


Update a integration:

```go
updatedIntegration := pingdomext.WebHookIntegration{
	Active:     true,
	ProviderID: 2,
	UserData: &pingdomext.WebHookData{
		Name: "tets-3",
		URL:  "http://www.example5.com",
	},
}
updateMsg, err := client_ext.Integrations.Update(12345, &updatedIntegration)
```

Delete a integration:

```go
delMsg, err := client_ext.Integrations.Delete(12345)
```

List all integration providers:

```go
listProviders, err := client_ext.Integrations.ListProviders()
```



### UserService ###

This service manages Solarwinds users. A Solarwinds user can be granted access to a range of services, each with its 
own web portal. Pingdom is one of those services. The Solarwinds API is a GraphQL API which is at a different domain 
as Pingdom API.

Create a new user invitation. The user will only appear on the active user list and be able to use Pingdom service after
he accepts the invitation manually.

```go
user := User{
    Email: "sombody@algolia.com",
    Role: "ADMIN"
    Products: []Product{
        {
        	Name: "PINGDOM",
        	Role: "MEMBER",
        }	
    }
}
err := client.UserService.Create(user)
```

Update an user. User information will be updated if the user has already accepted the invitation. If the invitation has
not yet been accepted, the invitation will be revoked and a new one with the updated information will be sent.
```go
update := User{
    Email: "sombody@algolia.com",
    Role: "ADMIN"
    Products: []Product{
    {
        Name: "PINGDOM",
        Role: "MEMBER",
    }
}
err := client.UserService.Update(update)
```

Delete an user. It is not possible to delete an active user in Solarwinds. If it is an active user, the function will
return error with proper status code set, no user will be deleted. If it is an invitation, the invitation will be revoked.

```go
email = "somebody@algolia.com"

err := client.UserService.Delete(email)
```

Retrieve an user. It can either be an invitation or an active user.

```go
email = "sombody@algolia.com"

err := client.UserService.Retrieve(email)
```

## Development ##

### Acceptance Tests ###

You can run acceptance tests against the actual pingdom API to test any changes:
```
PINGDOM_API_TOKEN=[api token] make acceptance
```

In order to run acceptance tests against the pingdom extension API, the following environment variables must be set:
```
SOLARWINDS_USER=[username] SOLARWINDS_PASSWD=[password] make acceptance
```

Note that this will create actual resources in your Pingdom account.  The tests will make a best effort to clean up but these would

In order to run acceptance tests against the actual Solarwinds API, the following environment variables must be set:
```
SOLARWINDS_USER=[solarwinds username] SOLARWINDS_PASSWD=[solarwinds password] make acceptance
```

Note that this will create actual resources in your Pingdom/Solarwinds account.  The tests will make a best effort to clean up but these would
not be guaranteed on test failures depending on the nature of the failure.
