package config

var SlackChannelIDs = map[string]string{
	"cloudplatform": "C07F60QARM3", // test-cloudplatform
	"k8s":           "C07F8PXH1CK", // test-k8s
	"prod":          "C07EUAWCFB9", // test-prod
	"winops":        "C07F8PW8LJX", // test-winops

	// Add more team name to channel ID mappings as needed
}

type Config struct {
	UseOllama     bool
	OpenAIKey     string
	GitHubRepoURL string
}

// WeeklyProgressTemplate is the template used for weekly progress updates in Markdown format
const WeeklyProgressTemplate = `
---
hide:
  - navigation
---
# Weekly progress updates

---

## Week of Monday 29th July 2024

Greetings, fellow Derivians! Wrapping up another week of Production Operations activities. Until next time, happy coding!

<div class="grid cards" markdown>

- 🚀 __Progress__

      ---

      ***<img src="/assets/images/winops.png" alt="Winops" width="20" height="20"> Team WinOps***

      - Integrated client latency reports for easier analysis.
      - Enhanced security by cleaning up default settings and adding network logs.
      - Addressed data feed piracy issues.

      ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Production***

      - Began migration preparations for key monitoring services.
      - Improved infrastructure with initial network setup changes.

      ***<img src="/assets/images/dba.png" alt="DBA" width="24" height="24"> Team DBA***

      - Optimized databases by cleaning up unused resources.
      - Completed maintenance on the VR database and resolved production issues.
      - Improved security score to 78% by fixing vulnerabilities.

      ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Development***

      - Migrated internal applications and improved stress testing processes.

      ***<img src="/assets/images/kubernetes.png" alt="Kubernetes" width="24" height="24"> Team PE Kubernetes Core***

      - Automated security updates.
      - Improved documentation for service migrations.
      ---

- ⚠️ __Problems__

      ---

      ***<img src="/assets/images/winops.png" alt="Winops" width="20" height="20"> Team WinOps***

      - Complex issues with unauthorized data access.
      - Potential cost increases due to server changes in the UAE.

      ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Production***

      - Reviewing feedback on security processes.

      ***<img src="/assets/images/dba.png" alt="DBA" width="24" height="24"> Team DBA***

      - Issues in pre-release testing for system updates.

      ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Development***

      - N/A
  
      ***<img src="/assets/images/kubernetes.png" alt="Kubernetes" width="24" height="24"> Team PE Kubernetes Core***

      - N/A

      ---

- 📋 __Plan__

      ---

      ***<img src="/assets/images/winops.png" alt="Winops" width="20" height="20"> Team WinOps***

      - Finalize the Real 03 requirements for relocation to UAE.
      - Enhance the monitoring tools in place to detect feed piracy cases.

      ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Production***

      - Finalize migration preparations and refine automation tools.

      ***<img src="/assets/images/dba.png" alt="DBA" width="24" height="24"> Team DBA***

      - Implement database improvements and finalize backup reviews.

      ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Development***

      - Complete stress testing planning for next steps.

      ***<img src="/assets/images/kubernetes.png" alt="Kubernetes" width="24" height="24"> Team PE Kubernetes Core***

      - Focus on documentation and support for service migrations.
      ---

-   💡 __Insights__

      ---

      ***<img src="/assets/images/winops.png" alt="Winops" width="20" height="20"> Team WinOps***

      - <span style="color: green;">**Highlights**</span>: Successful integration of latency reports (A Book and B Book)
      - <span style="color: red;">**Lowlights**</span>: Challenges with unauthorized piracy feed.

      ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Development***

      - <span style="color: green;">**Highlights**</span>: Progress in stress testing scenarios.
  

</div>



## <img src="/assets/images/blocker.png" alt="blocker" width="24" height="24"> **Challenges**
  🔴 Nothing blocks our way as of now 

---
`

// IndexMarkdown is the template used for the index page in Markdown format
const IndexMarkdown = `
---
hide:
 - navigation
toc_depth: 3
---

# Recent updates

Stay tuned as we share insights, challenges and successes in our quest for continuous improvement and innovation

???+ abstract "Latest monthly update"

    **July 2024**

    <div class="grid cards" markdown>

    - 🚀 __Progress [key initiatives, actions, achievements]__

        ---

        - **UAE License - MT5**: Finalizing UAE server infrastructure requirements. Primary servers are deployed, and relocation is scheduled for next week.

        - **MT5 Upgrade**: MT5 platforms are updated to the latest stable version **(4410)**.

        - **MT5 Mobile Signups**: For quite some time, the mobile signup process was misused for proprietary trading. Mobile signups  are now limited to 24 hours via automation, reducing server load by redirecting clients to the Deriv signup process for optimal load balancing. This reduced cost and optimization time on Demo servers.

        - **Kubernetes All Aboard**: The Platform Engineering team is charging full steam ahead with Kubernetes adoption and migration. This move will boost our security, agility, and developer experience.
            - Current Status: 
                - ***Internal Services***: Camunda, Redmine, and WikiJS are now running on Kubernetes. => Reduction of **$2.7K** in costs.
                - ***Production Services***: 16 services have been successfully migrated.
            - In Progress: 
                - Two external services, ***notification service*** and ***Hydra authentication*** tool, are being tested in QA and are gearing up for their big production debut. 

    - ⚠️ __Problems [blockers, risks, challenges]__

        ---

        - **CrowdStrike Incident**: The recent CrowdStrike incident showed both the strengths and weaknesses of our automation systems. While we managed to handle the situation quickly, we need to make our deployment pipeline more flexible and customizable for future scenarios. This incident highlights the risks of relying on third-party solutions and the need to build resilience to handle similar issues.

        - **Feed Piracy**: We're dealing with a piracy issue where synthetic feed is being stolen from our platforms by users connected as normal clients. To address this, the WinOps team is working with the Trading team to improve our monitoring processes. While some offenders have been identified, the investigation is ongoing to establish more comprehensive monitoring.

        - **MT5 Real accounts quota**: An issue with Real 01 platform's real account number quota arose after applying a business requirement to extend the archival period from 30 to 60 days; as a temporary solution, automation disables trading after 30 days of inactivity while keeping accounts unarchived for compliance purposes.


    - 📋 __Plan [key actions, objectives]__

        ---

        - **Latency and UAE License**: Failover of Real 03 to UAE to improve latency for EMEA Region  and offer UAE regulated accounts once the licence is obtained.

        - **South Africa Deployment - Real 03**: An additional trade server will balance the load among African servers and provide resilience during potential downtimes in our **highest trading activity region**.

        - **Optimization Time**: Add Demo Trade 04 to the load balancing to reduce optimization time for demo signups.
        - **Stress Testing and Fault Detection**: Finalise the stress testing plan for our infrastructure to gain a better understanding of its behaviour and limits. This will help us predict resource requirements and anticipate future activity growth.
        - **Kubernetes Migration**: We plan to continue migrating over 50 identified services to Kubernetes, aiming to have them running in our cluster by the end of August. 
        - **Feed DB Optimization**: We optimised the Tick Insert time from approximately 120ms to around 30 ms. This had a positive impact on our pricing. Our next goal is to reduce it further to below 30 ms and stabilise any outliers. 


    -   💡 __Insights [kpis, findings, trends]__

        ---

        ### KPIs
        - **Uptime for July 2024:**
            - MTS: 98.84% (Crowdstrike incident)
            - DTrader: 99.99%
        - **MT5 Optimization Time (max):**
            - Real: 12 min
            - Demo: 10 min
        - **MT5 Latency Time (A book, B book)**
            - Africa: 291ms (A), 273ms (B)
            - Americas: 216ms (A), 256ms (B)
            - Asia: 237ms (A), 284ms (B)
            - Europe: 124ms (A), 254ms (B)
        - **Migrated Kubernetes services:** 21

        ### Insights
        - **AWS Shield:** The solution offered by AWS is evaluated as an additional security layer for MT5 infrastructure.
        - **Chef:** We are removing our dependency to Chef which is one of our major security problems for years (#task_production_chef_zero)

    </div>






???+ abstract "Latest weekly update"

    **Week 2024-07-29 to 2024-08-02**

    <div class="grid cards" markdown>

    - 🚀 __Progress__

        ---

        ***<img src="/assets/images/winops.png" alt="Winops" width="20" height="20"> Team WinOps***

        - Integrated client latency reports for easier analysis.
        - Enhanced security by cleaning up default settings and adding network logs.
        - Addressed data feed piracy issues.

        ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Production***

        - Began migration preparations for key monitoring services.
        - Improved infrastructure with initial network setup changes.

        ***<img src="/assets/images/dba.png" alt="DBA" width="24" height="24"> Team DBA***

        - Optimized databases by cleaning up unused resources.
        - Completed maintenance on the VR database and resolved production issues.
        - Improved security score to 78% by fixing vulnerabilities.

        ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Development***

        - Migrated internal applications and improved stress testing processes.

        ***<img src="/assets/images/kubernetes.png" alt="Kubernetes" width="24" height="24"> Team PE Kubernetes Core***

        - Automated security updates.
        - Improved documentation for service migrations.
        ---

    - ⚠️ __Problems__

        ---

        ***<img src="/assets/images/winops.png" alt="Winops" width="20" height="20"> Team WinOps***

        - Complex issues with unauthorized data access.
        - Potential cost increases due to server changes in the UAE.

        ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Production***

        - Reviewing feedback on security processes.

        ***<img src="/assets/images/dba.png" alt="DBA" width="24" height="24"> Team DBA***

        - Issues in pre-release testing for system updates.

        ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Development***

        - N/A
    
        ***<img src="/assets/images/kubernetes.png" alt="Kubernetes" width="24" height="24"> Team PE Kubernetes Core***

        - N/A

        ---

    - 📋 __Plan__

        ---

        ***<img src="/assets/images/winops.png" alt="Winops" width="20" height="20"> Team WinOps***

        - Finalize the Real 03 requirements for relocation to UAE.
        - Enhance the monitoring tools in place to detect feed piracy cases.

        ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Production***

        - Finalize migration preparations and refine automation tools.

        ***<img src="/assets/images/dba.png" alt="DBA" width="24" height="24"> Team DBA***

        - Implement database improvements and finalize backup reviews.

        ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Development***

        - Complete stress testing planning for next steps.

        ***<img src="/assets/images/kubernetes.png" alt="Kubernetes" width="24" height="24"> Team PE Kubernetes Core***

        - Focus on documentation and support for service migrations.
        ---

    -   💡 __Insights__

        ---

        ***<img src="/assets/images/winops.png" alt="Winops" width="20" height="20"> Team WinOps***

        - <span style="color: green;">**Highlights**</span>: Successful integration of latency reports (A Book and B Book)
        - <span style="color: red;">**Lowlights**</span>: Challenges with unauthorized piracy feed.

        ***<img src="/assets/images/aws.png" alt="AWS" width="24" height="24"> Team PE Development***

        - <span style="color: green;">**Highlights**</span>: Progress in stress testing scenarios.
    

    </div>

## Previous updates

???+ abstract "August 2024"
    - ### [2024-07-029 to 2024-08-02](./2024-07-29.md)

???+ abstract "July 2024"
    - ### [2024-07-22 to 2024-07-26](./2024-07-22.md)
    - ### [2024-07-15 to 2024-07-19](./2024-07-15.md)
    - ### [2024-07-08 to 2024-07-12](./2024-07-08.md)
    - ### [2024-07-01 to 2024-07-05](./2024-07-01.md)

???+ abstract "June 2024"
    - ### [2024-06-24 to 2024-06-28](./2024-06-28.md)
    - ### [2024-06-17 to 2024-06-21](./2024-06-17.md)
    - ### [2024-06-10 to 2024-06-14](./2024-06-10.md)
    - ### [2024-06-03 to 2024-06-07](./2024-06-03.md)


???+ abstract "May 2024"
    - ### [2024-05-27 to 2024-05-31](./2024-05-27.md)
    - ### [2024-05-20 to 2024-05-24](./2024-05-20.md)
    - ### [2024-05-13 to 2024-05-17](./2024-05-13.md)
    - ### [2024-05-06 to 2024-05-10](./2024-05-06.md)


???+ abstract "April 2024"
    - ### [2024-04-29 to 2024-05-03](./2024-04-29.md)
    - ### [2024-04-22 to 2024-04-26](./2024-04-22.md)
    - ### [2024-04-15 to 2024-04-19](./2024-04-15.md)
`
