# Research how indexes affect performance
## Task

Generate 1,000,000 profiles by any means. First and Last names must be real to account for index selectivity. You can also use an existing list as a basis.

Implement functionality to search profiles by first name and last name prefix (simultaneously) in your social network (implement the `/user/search` method from the specification) (query in the form `firstName LIKE ? and secondName LIKE ?`). Sort the output by profile id.

Conduct load tests for this method. Experiment with the number of concurrent requests: 1/10/100/1000.

Build graphs and save them in the report.

Create a suitable index.

Repeat steps 3 and 4.

As a result, provide a report that must include:

- Latency graphs before the index;
- Throughput graphs before the index;
- Latency graphs after the index;
- Throughput graphs after the index;
- The index creation query;
- Explain queries after the index;
- Explanation of why this particular index was chosen.

## Solution
