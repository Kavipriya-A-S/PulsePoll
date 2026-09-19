import { useEffect, useState } from "react";
import "./App.css";

// Deployed Go/Gin backend
const API_URL = "https://pulsepoll-syhz.onrender.com/api";

function App() {
  const [polls, setPolls] = useState([]);
  const [question, setQuestion] = useState("");
  const [options, setOptions] = useState(["", ""]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  // Fetch all polls
  const fetchPolls = async () => {
    try {
      setError("");

      const response = await fetch(`${API_URL}/polls`);

      if (!response.ok) {
        throw new Error("Failed to fetch polls");
      }

      const data = await response.json();
      setPolls(data);
    } catch (err) {
      console.error("Error fetching polls:", err);
      setError("Unable to connect to the PulsePoll backend.");
    }
  };

  // Fetch polls when page loads
  useEffect(() => {
    const loadPolls = async () => {
      try {
        setError("");

        const response = await fetch(`${API_URL}/polls`);

        if (!response.ok) {
          throw new Error("Failed to fetch polls");
        }

        const data = await response.json();
        setPolls(data);
      } catch (err) {
        console.error("Error fetching polls:", err);
        setError("Unable to connect to the PulsePoll backend.");
      }
    };

    loadPolls();
  }, []);

  // Update an option
  const handleOptionChange = (index, value) => {
    const updatedOptions = [...options];
    updatedOptions[index] = value;
    setOptions(updatedOptions);
  };

  // Add a new option
  const addOption = () => {
    setOptions([...options, ""]);
  };

  // Remove an option
  const removeOption = (index) => {
    if (options.length <= 2) return;

    setOptions(options.filter((_, i) => i !== index));
  };

  // Create a new poll
  const createPoll = async (e) => {
    e.preventDefault();

    setMessage("");
    setError("");

    const cleanedOptions = options
      .map((option) => option.trim())
      .filter(Boolean);

    if (!question.trim()) {
      setError("Please enter a poll question.");
      return;
    }

    if (cleanedOptions.length < 2) {
      setError("Please provide at least two options.");
      return;
    }

    try {
      setLoading(true);

      const response = await fetch(`${API_URL}/polls`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          question: question.trim(),
          options: cleanedOptions.map((text) => ({
            text,
          })),
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Failed to create poll");
      }

      setQuestion("");
      setOptions(["", ""]);
      setMessage("Poll created successfully!");

      await fetchPolls();
    } catch (err) {
      console.error("Error creating poll:", err);
      setError(err.message || "Failed to create poll.");
    } finally {
      setLoading(false);
    }
  };

  // Vote on a poll option
  const vote = async (pollId, optionId) => {
    try {
      setError("");
      setMessage("");

      const response = await fetch(
        `${API_URL}/polls/${pollId}/vote/${optionId}`,
        {
          method: "POST",
        }
      );

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Failed to record vote");
      }

      setMessage("Your vote has been recorded!");

      await fetchPolls();
    } catch (err) {
      console.error("Error recording vote:", err);
      setError(err.message || "Failed to record vote.");
    }
  };

  // Calculate total votes
  const getTotalVotes = (poll) => {
    return poll.options.reduce(
      (total, option) => total + option.votes,
      0
    );
  };

  return (
    <div className="app">
      {/* Navbar */}
      <header className="navbar">
        <div className="logo">
          Pulse<span>Poll</span>
        </div>

        <div className="nav-status">
          <span className="status-dot"></span>
          Live Polling
        </div>
      </header>

      <main className="container">
        {/* Hero */}
        <section className="hero">
          <p className="eyebrow">REAL-TIME POLLING PLATFORM</p>

          <h1>
            Ask. Vote.
            <br />
            <span>See the Pulse.</span>
          </h1>

          <p className="hero-text">
            Create live polls, collect votes, and instantly see what people
            think.
          </p>
        </section>

        {/* Messages */}
        {message && <div className="success-message">{message}</div>}

        {error && <div className="error-message">{error}</div>}

        {/* Create Poll */}
        <section className="create-section">
          <div className="section-heading">
            <div>
              <p className="section-label">CREATE</p>
              <h2>Start a new poll</h2>
            </div>
          </div>

          <form className="poll-form" onSubmit={createPoll}>
            {/* Question */}
            <label>
              Poll question

              <input
                type="text"
                placeholder="What would you like to ask?"
                value={question}
                onChange={(e) => setQuestion(e.target.value)}
              />
            </label>

            {/* Options Header */}
            <div className="options-header">
              <span>Options</span>

              <button
                type="button"
                className="add-option"
                onClick={addOption}
              >
                + Add option
              </button>
            </div>

            {/* Options */}
            <div className="option-list">
              {options.map((option, index) => (
                <div className="option-input" key={index}>
                  <span>{index + 1}</span>

                  <input
                    type="text"
                    placeholder={`Option ${index + 1}`}
                    value={option}
                    onChange={(e) =>
                      handleOptionChange(index, e.target.value)
                    }
                  />

                  {options.length > 2 && (
                    <button
                      type="button"
                      className="remove-option"
                      onClick={() => removeOption(index)}
                    >
                      ×
                    </button>
                  )}
                </div>
              ))}
            </div>

            {/* Create Button */}
            <button
              className="create-button"
              type="submit"
              disabled={loading}
            >
              {loading ? "Creating..." : "Create Poll"}
            </button>
          </form>
        </section>

        {/* Live Polls */}
        <section className="polls-section">
          <div className="section-heading polls-heading">
            <div>
              <p className="section-label">EXPLORE</p>
              <h2>Live polls</h2>
            </div>

            <button
              type="button"
              className="refresh-button"
              onClick={fetchPolls}
            >
              ↻ Refresh
            </button>
          </div>

          {/* No Polls */}
          {polls.length === 0 ? (
            <div className="empty-state">
              <div className="empty-icon">◉</div>

              <h3>No polls yet</h3>

              <p>
                Create your first poll and start collecting votes.
              </p>
            </div>
          ) : (
            /* Poll Grid */
            <div className="poll-grid">
              {polls.map((poll) => {
                const totalVotes = getTotalVotes(poll);

                return (
                  <article className="poll-card" key={poll.id}>
                    {/* Poll Header */}
                    <div className="poll-card-top">
                      <span className="live-badge">
                        <span></span> LIVE
                      </span>

                      <span className="vote-count">
                        {totalVotes}{" "}
                        {totalVotes === 1 ? "vote" : "votes"}
                      </span>
                    </div>

                    {/* Question */}
                    <h3>{poll.question}</h3>

                    {/* Voting Options */}
                    <div className="vote-options">
                      {poll.options.map((option) => {
                        const percentage =
                          totalVotes > 0
                            ? Math.round(
                                (option.votes / totalVotes) * 100
                              )
                            : 0;

                        return (
                          <div
                            className="vote-option"
                            key={option.id}
                          >
                            <button
                              type="button"
                              onClick={() =>
                                vote(poll.id, option.id)
                              }
                            >
                              <span>{option.text}</span>

                              <strong>{percentage}%</strong>
                            </button>

                            <div className="progress-track">
                              <div
                                className="progress-bar"
                                style={{
                                  width: `${percentage}%`,
                                }}
                              ></div>
                            </div>
                          </div>
                        );
                      })}
                    </div>

                    {/* Footer */}
                    <div className="poll-footer">
                      Select an option to vote
                    </div>
                  </article>
                );
              })}
            </div>
          )}
        </section>
      </main>

      {/* Footer */}
      <footer>
        <p>PulsePoll · Live opinions, one vote at a time.</p>
      </footer>
    </div>
  );
}

export default App;