// Domain model

job
  (id)
  task_form

  verification_form
  verification_module
  verification_agreement_count
  verification_score_threshold
  funding_verification_reward
  verification_limit

  task
    (job_id, id)
    data
    is_active
    assignment_count
    pending_count
    accepted_count
    verification_count

//- task flow
assignment
  (job_id, task_id, id)
  user_id

responses
  (job_id, task_id, id)
  worker_id
  value
  is_accepted
  verifications_count

//- verification flow
verification_assignment
  (job_id, task_id, response_id, id)
  user_id  

verification_responses
  (job_id, task_id, response_id, id)
  worker_id
  verifier_id
  value
  reason

verificaiton_eligibility
  job_id
  user_id


// modules --------------------------------
verification-settings
  get(jobId)
  create(settings)
  
  settings: 
    type (automatic (consenus, external), manual)

    manual
      whitelist?
      form
      verification_limit

    automatic
      agreement_count



verification-assignment
  assign(jobId, userId)
  unassign(jobId, userId)

verification-verify
  verify(userId, responseId, result, reason)
  verifyAtomatic()


// current
[POST] /api/v1/tasks/:taskID/submit
  // if verificationsettings().automatic
        verificationModule.Verify() // automatic

[POST] /api/v1/responses/:responseID/verify
  // Ensure verification is assigned

  verificationModule.Verify() // manual
    // if need payout ? 
    //   payoutVerifier
    scoreResponse()
      ScoreResponse()
        // Insert score
        // Set response to accepted
        // Delete verification assignment
        // Update task counters (pending)
        // Update user counters (accepted, rejected) 
        // Update Assignment (svc)

      // payout or reject task (funding)
      // fireEvent(accepted)


// new
[POST] /api/v1/tasks/:taskID/submit
  // if !verificationsvc.settings().manual
        verificationsvc.VerifyAutomatic()

[POST] /api/v1/responses/:responseID/verify
  // Ensure verification is assigned

  settings = verificationsvc.settings()

  if settings.manual
    if !settings.requester ? 
      payoutVerifier()

    verificationsvc.SaveScore()

    // Set response to accepted
    // Delete verification assignment
    // Update task counters (pending)
    // Update user counters (accepted, rejected) 
    // Update Assignment (svc)

    // payout or reject task (funding)
    // fireEvent(accepted)

